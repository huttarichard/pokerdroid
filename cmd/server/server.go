package cmdserver

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/bot/cfr"
	"github.com/pokerdroid/poker/bot/mc"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/tree"
	"github.com/spf13/cobra"
)

type serverArgs struct {
	abs  string
	dir  string
	addr string
}

var tf = serverArgs{}

func init() {
	flags := serverCMD.Flags()

	flags.StringVar(&tf.abs, "abs", "", "path to the abstraction")
	flags.StringVar(&tf.dir, "dir", "", "path to the directory with solutions")

	flags.StringVar(&tf.addr, "addr", ":8080", "address to listen on")

	cobra.MarkFlagRequired(flags, "abs")
	cobra.MarkFlagRequired(flags, "dir")

}

var serverCMD = &cobra.Command{
	Use:   "serve",
	Short: "will run server that will serve solutions over HTTP",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		logger := log.Default()

		logger.Printf("loading abstraction")
		abs, err := absp.NewFromFile(tf.abs)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("abs: %s", abs.UID.String())

		rng := frand.NewHash()

		rxs, err := tree.NewFileRootsFromDir(tf.dir)
		if err != nil {
			logger.Fatal(err)

		}
		defer rxs.Close()

		for _, rx := range rxs.Roots() {
			logger.Printf("loaded root from %s", rx.Params.String())
			logger.Printf("absid: %s", rx.AbsID.String())
		}

		cfradv := cfr.Advisor{
			Roots:   rxs.Roots(),
			Abs:     abs,
			Rand:    rng,
			Logger:  logger,
			Advisor: cfr.AdvisorSimple,
		}

		mcadv := mc.NewAdvisor()

		combined := bot.NewCombined(cfradv, mcadv)

		mux := bot.NewHTTP(combined)
		mux.Logger = logger

		r := chi.NewRouter()
		r.Mount("/advise", mux)

		server := &http.Server{
			Addr:    tf.addr,
			Handler: r,
		}

		logger.Printf("starting server on %s", tf.addr)
		go server.ListenAndServe()

		<-ctx.Done()
		server.Shutdown(ctx)
	},
}
