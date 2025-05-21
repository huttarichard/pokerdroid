package abs

import "github.com/pokerdroid/poker/card"

type Cluster uint32

type Clusters []Cluster

type Mapper interface {
	Map(cds card.Cards) Cluster
}
