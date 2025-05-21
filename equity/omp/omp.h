#ifndef OMP_H
#define OMP_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

double hand_equity(const char *hole_card, const char *community_card, const char *dead_card, const int nb_players,
		const int nb_board_cards_, const double std_tol);

void hand_equity_with_draw(const char *hole_card, const char *community_card, const char *dead_card, 
                          const int nb_players, double *win_equity, double *draw_equity,
                          const int nb_board_cards_, const double std_tol);

#ifdef __cplusplus
}
#endif

#endif  /* OMP_H */
