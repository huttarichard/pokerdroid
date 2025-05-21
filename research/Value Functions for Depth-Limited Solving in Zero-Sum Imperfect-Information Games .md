# Value Functions for Depth-Limited Solving in Zero-Sum Imperfect-Information Games 

Vojtěch Kovařík ${ }^{1}$, Dominik Seitz ${ }^{1}$, Viliam Lisý*, Jan Rudolf, Shuo Sun, Karel Ha<br>Artificial Intelligence Center, FEE, Czech Technical University in Prague, Prague, Czech Republic


#### Abstract

We provide a formal definition of depth-limited games together with an accessible and rigorous explanation of the underlying concepts, both of which were previously missing in imperfect-information games. The definition works for an arbitrary extensive-form game and is not tied to any specific game-solving algorithm. Moreover, this framework unifies and significantly extends three approaches to depth-limited solving that previously existed in extensive-form games and multiagent reinforcement learning but were not known to be compatible. A key ingredient of these depth-limited games are value functions. Focusing on two-player zero-sum imperfect-information games, we show how to obtain optimal value functions and prove that public information provides both necessary and sufficient context for computing them. We provide a domain-independent encoding of the domains that allows for approximating value functions even by simple feed-forward neural networks, which are then able to generalize to unseen parts of the game. We use the resulting value network to implement a depthlimited version of counterfactual regret minimization. In three distinct domains, we show that the algorithm's exploitability is roughly linearly dependent on the value network's quality and that it is not difficult to train a value network with which depth-limited CFR's performance is as good as that of CFR with access to the full game.


Keywords: Imperfect Information Game, Multiagent Reinforcement Learning, Extensive Form Game, Partially Observable Stochastic Game, Depth Limited Game, Depth Limited Solving, Value Function, Counterfactual Regret Minimization

[^0]
## 1. Introduction

Sequential decision making is a key challenge in AI research. As the number of consequent decisions increases, the size of the state space blows up exponentially, to the point where even modern computer clusters soon become unable to even enumerate all the states. In perfect information problems, this issue is typically overcome by replacing the states below a certain depth by a value function. This technique vastly reduces the effective size of the game, which is essential for both minimax-like search and reinforcement learning. In imperfect information problems, value functions are much more complex since they depend on the agent's belief about the current state. The situation gets even more complicated in multiagent imperfect information problems, where values additionally depend on each agent's belief about other agents' beliefs (etc.). Despite these challenges, recent results in poker [1, 2, 3] illustrate that depth-limited approaches can be successful even in this setting.

Unfortunately, many of the key concepts that enabled the recent results were introduced informally or tied to specific domains. As a result, it is unclear how to adapt these results to new domains, the process is time-consuming, and the theoretical properties of the resulting algorithms are unknown. For example, the value function used in [1] takes as input the probabilities of different poker hands that a player could be holding. This makes sense in poker, but what would be the input if we tried a similar approach in blind chess, scrabble, or phantom tic-tac-toe? How would we train a value function in these domains and how would we know whether it is "good"? And what if we wanted to use the value function in combination with a different algorithm? How would we do it, and would it "work"? What does "good" and "work" mean in this context? In this work, we aim to answer these questions by providing a solid theory of depth limited games and value functions and illustrating it by experiments on several domains.

### 1.1. Outline and Contributions

In summary, this paper (1) provides a theory of depth-limited methods and value functions that unifies three recent approaches [1, 4, 5], (2) formulates all required concepts in an accessible and domain-independent way, and (3) experimentally demonstrates that depth-limited solving is a viable and robust option for a range of imperfect information games beyond poker, and can be done without the need for hand-crafting domain-specific features.

In more detail, we start with a brief background on EFGs (Section 2). In Section 3] we define the key concepts - expected utilities of histories and information sets, reach probabilities, and beliefs - in a unified way that is consistent with previous literature but much more intuitive and easy to use. We then provide several technical propositions which capture the intuitive properties of values in EFGs and substantially simplify our proofs. A more detailed outline of this section is provided in its introduction (and the same is true for the other longer sections, i.e., Sections 4 and Section 6).

In Section 4, we present the key theoretical contributions of the paper. First, we look at imperfect-information games and propose domain- and algorithmindependent notions of value functions and depth-limited games. We also describe (the depth-limited versions of) various algorithmic problems such as game-value computation, equilibrium computation, and best response computation. Our goal for each of these problems is that whenever we find a solution of the depth-limited version of the problem and plug it into the full game, it should fulfil the role of a (partial) solution to the non-depth-limited version of the problem. For example, a depth-limited Nash equilibrium should coincide with a standard Nash equilibrium in all decision-points above the depth limit. However, even some promising choices of value functions can fail to work for some of the above problems. We thus describe a natural hierarchy of conditions on value functions and formally prove that each of them fulfils the above goal for a different class of computational problems. Moreover, we observe that with the proposed formalization, the two previously separate approaches to depth-limited solving of imperfect information games - value functions [1] and multivalued states [4] can be viewed as two instances of a single unifying framework. After studying the fundamental properties of value functions, we discuss methods for representing value functions more efficiently. One part of this endeavour is encoding the functions' input and output more compactly - we show that value functions can be defined on either individual histories or information sets and that the two representations can be translated to each other. Moreover, we formally prove that value functions can be factorized based on public information (or common knowledge) and that this factorization cannot be further refined in general. The second part is approximating value functions by neural networks, which is likely to be easier if there is a unique approximation target. Unfortunately, we see that there can sometimes be multiple suitable value functions. Failing uniqueness, we investigate whether the set of suitable value functions is at least well-behaved i.e., convex. We prove that this is true for certain types of value functions and pose the general question as an open problem. However, we remain optimistic about value function approximation since the non-uniqueness did not prove to be a problem in practice.

While all of the above results are presented using the extensive-form game formalism, Section 5 explains how all of them apply the POSG- (partiallyobservable stochastic game; [6]) and FOSG- (factored-observation stochastic game; [7) formalisms used in multiagent reinforcement learning.

The experimental part of the paper (Section 6) focuses on the depth-limited version of counterfactual regret minimization (CFR), an algorithm which seems particularly promising due to its recent successes in poker [1, 2, 8] and other domains [9, 10]. We show that any game represented as a FOSG admits a unified representation of inputs and outputs to the value function. We demonstrate, in three domains with very different properties, that this representation can serve as a suitable encoding for a neural network. This encoding allows for an accurate approximation even with the simplest feed-forward architecture. We show that when this value network is plugged into depth-limited CFR, the algorithm produces strategies with low exploitability. While doing so, we
experiment with various loss functions in an attempt to identify those which serve best as a proxy for exploitability when used in conjunction with depthlimited CFR. To our surprise, the algorithm's performance is not very sensitive to this parameter. We also investigate whether the trained value function is robust to differences between the training distribution and the one demanded by depth limited CFR. In all three domains, the trained value function generalizes well to unseen situations. Furthermore, we investigate the dependence of the quality of the learned value function on the amount of training data and key hyper-parameters.

Finally, we summarize the most-related existing works (Section 7.1) and present our conclusions (Section 7.2).

### 1.2. Novelty of the Theoretical Contribution

Since the main contribution of Section 3 is in formalizing and organizing content that is well-known or simple, the results presented there are not novel unless stated otherwise. The corresponding novel content is listed in the introduction to this section. On the other hand, all content in Section 4 (and Appendix Ap) is novel unless stated otherwise. Finally, an important novel part of the paper is the domain-independent encoding proposed in Section 6.1.5.

## 2. Extensive-Form Games

In this section, we formally describe the EFG model used throughout the paper. We make a slight deviation from the historically standard (but much less convenient) definition by assuming that information sets of each player are defined even over terminal states and states in which it is the opponent's turn. This modification is consistent with (4] and [7.

A extensive-form game $G$ can be described by: $\mathcal{H}$ - a finite set of histories (states, nodes), representing sequences of actions. We use $g \sqsubset h$ to denote the fact that $g$ is equal to or a prefix of $h . \mathcal{Z}$ - the set of terminal histories (those $z \in \mathcal{H}$ which are not a prefix of any other history). $\mathcal{N}=\{1, \ldots, N\}-$ the player set, where $c$ is a special player, called "chance" or "nature". The player function $\mathcal{P}: \mathcal{H} \backslash \mathcal{Z} \rightarrow \mathcal{N} \cup\{c\}$ denotes which player acts in the given non-terminal history. $\mathcal{A}(h):=\mathcal{A}_{p}(h):=\{a \mid h a \in \mathcal{H}\}$ is the set of actions available to $p=\mathcal{P}(h)$ at $h h^{2}$ The strategy of chance is a fixed probability distribution $\sigma_{c}$ over actions available at chance nodes (those $h$ where $\mathcal{P}(h)=c$ ). The utility function $u=\left(u_{p}\right)_{p \in \mathcal{N}}$ assigns to each terminal history $z$ a reward $u_{p}(z) \in \mathbb{R}$ received by each player upon reaching $z$.

The information-partitions $\mathcal{I}=\left(\mathcal{I}_{p}\right)_{p \in \mathcal{N}}$, where each $\mathcal{I}_{p}$ is a partition of $\mathcal{H}$, capture the imperfect information of $G$. If $g, h \in \mathcal{H}$ belong to the same information set (infoset) $I \in \mathcal{I}_{p}$ then $p$ cannot distinguish between them. For

[^1]each $I \in \mathcal{I}_{p}$, the available actions $\mathcal{A}_{p}(I):=\mathcal{A}_{p}(h)$ are thus the same for each $h \in I$. To help decompose the game, we consider the public tree $\mathcal{S}$ which partitions $\mathcal{H}$ into public states $s^{3} S \in \mathcal{S}$ which are closed ${ }^{4}$ under membership within infosets. We only consider games where players have perfect recall, i.e., they remember their past actions and infosets visited so far 5 A behavioral strategy $\sigma_{p} \in \Sigma_{p}$ of player $p$ assigns to each $I \in \mathcal{I}_{p}$ a probability distribution $\sigma_{p}(I)$ over actions $\mathcal{A}_{p}(I)$. A tuple $\sigma=\left(\sigma_{p}\right)_{p \in \mathcal{N}}$ is called a strategy profile. By $u_{p}(\sigma)$, we denote the expected utility for player $p$ if all players play according to $\sigma$.

The profile $\sigma$ is an $\boldsymbol{\epsilon}$-Nash equilibrium if the benefit of switching to some alternative $\sigma_{p}^{\prime}$ is limited by $\epsilon$, i.e., if $u_{p}\left(\sigma_{p}^{\prime}, \sigma_{-p}\right) \leq u_{p}(\sigma)+\epsilon$ (where $-p$ denotes all the players other than $p$ ). When $\epsilon=0$, the profile is called a Nash equilibrium (NE) and we write $\sigma \in \mathrm{NE}(G)$. In the remainder of the paper, we only considert two-player zero-sum games, where $\mathcal{N}=\{1,2\}$ and $u_{2}=-u_{1}$. It is a standard result that in two-player zero-sum games, all $\sigma^{*} \in \mathrm{NE}(G)$ have the same utility, called game value $u_{p}\left(\sigma^{*}\right)=\max _{\sigma_{p} \in \Sigma_{p}} \min _{\sigma_{-p} \in \Sigma_{-p}} u_{p}(\sigma)=: \operatorname{gv}_{p}(G)$. The exploitability $\operatorname{expl}(\sigma)$ of $\sigma \in \Sigma$ is the average of $\operatorname{exploitabilities~} \operatorname{expl}_{p}\left(\sigma_{p}\right)$, where $\operatorname{expl} l_{p}\left(\sigma_{p}\right):=\operatorname{gv}_{p}(G)-\min _{\sigma_{-p} \in \Sigma_{-p}} u_{p}\left(\sigma_{p}, \sigma_{-p}\right)$.

## 3. Clarifying Key Concepts in Imperfect Information Games

In this section, we describe the remaining prerequisites of this paper. Unfortunately, while most of them are considered standard and well-known by the community around EFGs and CFR, they have often not been previously published or made fully formal, and some of the existing definitions cannot be applied to unreachable parts of the game (which prevents reasoning about counterfactual scenarios where one of the players changes their strategy). Moreover, many of the concepts are "not very accessible" to new audience, largely because their definitions are scattered across many different conference papers with strict page limits. The goal of this section is to remedy this situation by explaining all of these concepts in one place, in a way that would be fully formal and accessible to readers not yet fully familiar with the CFR literature.

Given the goal of this section, it is unsurprising that most of the content is not novel - our contribution is in organizing it, finding the right formulations (which was sometimes surprisingly difficult), and coming up with the formal proofs (which was typically straightforward once we had the right formulations).

[^2]However, to our best knowledge, several of the observations are novel, and might be of interest even to readers well-versed in the CFR literature ${ }^{7}$ These are: (a) The observation that a number of properties of historyand infoset- values are not just a corollary of the definition, but in fact equivalent to it (the "moreover" parts of Lemma 3.1 and Theorem 11. (b) The observation that $p$ 's reach probability of (info)sets should be defined as in eq. 3.19. not as a sum of $p$ 's reach probabilities of histories in that (info)set. (c) Using an approach inspired by trembling-hand equilibria to extend beliefs to unreachable infosets (Definition 3.4). (d) The observation that once we know counterfactual values of leaf infosets, backpropagating them through the infoset tree is simple (and much simpler than backpropagating normal values; Theorem 2 .

We first describe the key concepts for histories (Section 3.1). We then extend these to sets of histories, and in particular to infosets (Section 3.2). In Section 3.3 , we describe counterfactual values and the counterfactual regret minimization algorithm, which plays an important role in subsequent sections.

### 3.1. Expected Utilities of Histories

Before describing expected utilities, we need the concept of a reach probability of a history. This definition and others in this section rely on the intuition that we are playing the game $G$ using some strategy (in this case $\sigma$ ), we are currently in some situation in $G$ (here the root), and we are asking about the probability of some event in this playthrough of $G$ (here encountering $h$ ). To avoid confusion, the less-formal definition (here eq. 3.1) is always accompanied by its fully-formal equivalent (here eq. 3.2 ):

$$
\begin{align*}
P^{\sigma}(h) & :=\mathbf{P r}_{\sigma}[h \text { reached during the course of the game }]  \tag{3.1}\\
& :=\Pi_{g a \sqsubset h} \sigma(g, a) \square^{8} \tag{3.2}
\end{align*}
$$

Reach probabilities can be decomposed into the player- $p$ component, and the corresponding $p$ 's counterfactual reach probability (i.e., the chance of reaching $h$ in the counterfactual (cf.) situation where $p$ aims to do so):

$$
\begin{align*}
P_{p}^{\sigma}(h) & :=\mathbf{P r}_{\sigma}\left[h \text { reached } \mid \text { every } p^{\prime} \neq p \text { plays to reach } h\right]  \tag{3.3}\\
& :=\Pi_{g a \sqsubset h, \mathcal{P}(g)=p} \sigma(g, a),  \tag{3.4}\\
P_{-p}^{\sigma}(h) & :=\mathbf{P r}_{\sigma}[h \text { reached } \mid p \text { plays to reach } h]  \tag{3.5}\\
& :=\Pi_{g a \sqsubset h, \mathcal{P}(g) \neq p} \sigma(g, a) . \tag{3.6}
\end{align*}
$$

These definitions implies that for every $p$,

$$
\begin{equation*}
P^{\sigma}(h)=P_{1}^{\sigma}(h) P_{2}^{\sigma}(h) P_{c}^{\sigma}(h)=P_{p}^{\sigma}(h) P_{-p}^{\sigma}(h) \tag{3.7}
\end{equation*}
$$

[^3]We say that a node $h$ is reachable, resp. counterfactually reachable by $p$, under $\sigma$ if $P^{\sigma}(h)>0$, resp. $P_{-p}^{\sigma}(h)>0$. When the reach probability is zero, $h$ is said to be unreachable, resp. counterfactually unreachable. We extend reach probabilities to paths ${ }^{9}$, obtaining an analogue of (3.7):

$$
\begin{align*}
& P^{\sigma}(g, h):=\mathbf{P r}_{\sigma}[h \text { reached } \mid \text { current node is } g]=\Pi_{g \sqsubset h^{\prime}, h^{\prime} a \sqsubset h} \sigma\left(h^{\prime}, a\right),  \tag{3.8}\\
& P^{\sigma}(g, h)=P_{1}^{\sigma}(g, h) P_{2}^{\sigma}(g, h) P_{c}^{\sigma}(g, h)=P_{p}^{\sigma}(g, h) P_{-p}^{\sigma}(g, h) . \tag{3.9}
\end{align*}
$$

As a useful piece of terminology, we say that a set $H \subset \mathcal{H}$ is thin if no two nodes $g, h \in H, g \neq h$, satisfy $g \sqsubset h$ (i.e., $H$ is an antichain w.r.t. $\sqsubset)$. Otherwise, $H$ is thick. A thin set $L \subset \mathcal{H}$ to which no $h \notin L$ can be added without making it thick (i.e., a maximal antichain w.r.t. $\sqsubset)$ is called a slice (through $\mathcal{H}$ ). While we typically imagine that $L$ slices $\mathcal{H}$ "somewhere in the middle", the trivial examples of slices are the singleton $\{$ root $\}$ and the set of all leaves $\mathcal{Z}$. However, slices do not need to be strictly "horizontal" - that is, they can contain histories with different lengths.

We now define the value of $h \in \mathcal{H}$ (for $p \in \mathcal{N}$ under $\sigma$ ) and the corresponding action values (when $h \notin \mathcal{Z}$ ) as $s^{10}$

$$
\begin{align*}
v_{p}^{\sigma}(h) & :=\mathbf{E}_{\sigma}\left[u_{p}(z) \mid z \in \mathcal{Z}, \text { current history is } h\right]  \tag{3.10}\\
& :=\sum_{z \in \mathcal{Z}} P^{\sigma}(h, z) u_{p}(z),  \tag{3.11}\\
q_{p}^{\sigma}(h, a) & :=\mathbf{E}_{\sigma}\left[u_{p}(z) \mid z \in \mathcal{Z}, \text { current history is } h, a \text { taken at } h\right]  \tag{3.12}\\
& :=\sum_{z \in \mathcal{Z}} P^{\sigma}(h a, z) u_{p}(z) . \tag{3.13}
\end{align*}
$$

Lemma 3.1 summarizes the properties of expected values of histories: (1) states that, by definition, values are calculated as expectations over utilities of terminal states. Moreover, by (2), this is equivalent to each history-value being the expectation of values over the history's children. (2') restates (2) from $p$ 's point of view, by saying that these expectations can be expressed as weighted sums, where the weights are either $p$ 's action probabilities or the counterfactual reach probabilities probabilities of reaching the child from $h$. An important consequence of these properties is (3), which claims that values can also be computed as expectations over any slice $L$ through the game tree (below the given history $h$, i.e., satisfying $\nexists g \in L \backslash\{h\}: g \sqsubset h)$.

Lemma 3.1 (Characterization of $v_{p}^{\sigma}$ ). For any $p \in \mathcal{N}$ and $\sigma \in \Sigma$, $v$ (root) $=$ $u_{p}(\sigma)$ and the values $v_{p}^{\sigma}$ and $q_{p}^{\sigma}(h, a):=v_{p}^{\sigma}(h a)$ have the following properties

[^4](1) $v_{p}^{\sigma}(h)=\mathbf{E}_{\sigma}\left[u_{p}(z) \mid z \in \mathcal{Z}\right.$, current history is $\left.h\right]$ for every $h \in \mathcal{H}$.
(2) $v_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and on $\mathcal{H} \backslash \mathcal{Z}$, we have $v_{p}^{\sigma}(h)=\sum_{a \in \mathcal{A}(h)} \sigma(h, a) v_{p}^{\sigma}(h a)$.
(2') $v_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and on $\mathcal{H} \backslash \mathcal{Z}$, we have
\$\$

v_{p}^{\sigma}(h)= $$
\begin{cases}\sum_{a \in \mathcal{A}_{p}(h)} \sigma_{p}(h, a) q_{p}^{\sigma}(h, a) & \text { when } \mathcal{P}(h)=p \\ \sum_{a \in \mathcal{A}(h)} P_{-p}^{\sigma}(h, h a) v_{p}^{\sigma}(h a) & \text { when } \mathcal{P}(h) \neq p\end{cases}
$$

\$\$
(3) $v_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and on $\mathcal{H} \backslash \mathcal{Z}$, we have $v_{p}^{\sigma}(h)=\sum_{h^{\prime} \in L} P^{\sigma}\left(h, h^{\prime}\right) v_{p}^{\sigma}\left(h^{\prime}\right)$ for every slice $L$ through $\mathcal{H}$ below $h$.

Moreover, each of these conditions can be used as an equivalent definition of $v_{p}^{\sigma}$ (i.e., it automatically implies all the others).

These properties and their equivalence are used - typically implicitly - in essentially all related proofs (e.g., they imply that it is correct to talk about the value of a history $h$ even when the strategy above $h$ is unknown). While this result is by no means surprising, it provides an intuition for what properties we might get for infosets values.

Proof. First, applying (1) to $h=$ root gives $v($ root $)=u_{p}(\sigma)$. Second, note that (2) and (2') are equivalent (since $\sigma(h, a)$ is equal to either $\sigma_{p}(h, a)$ or $P_{-p}^{\sigma}(h, h a)$, depending on $\mathcal{P}(h)$ and $v_{p}^{\sigma}(h a)=q_{p}^{\sigma}(h, a)$ when $\left.\mathcal{P}(h)\right)$.

Finally, for the remaining equivalences, consider the following lemma (whose proof is in the appendix, together with all other proofs not shown in the main text).

Lemma 3.2. Let $(T, \sqsubset)$ be a finite tree, $Z \subset T$ its leaves, $f: Z \rightarrow \mathbb{R}$, and $P: T^{2} \rightarrow[0,1]$ a function s.t. $P(t, t)=1, P(t, s)=0$ when $\neg(t \sqsubset s)$, and $P(s, u)=P(s, t) P(t, u)$ when $s \sqsubset t \sqsubset u$. Then the following are equivalent for $F: T \rightarrow \mathbb{R}$ :
(a) $F(t)=\sum_{z \in Z} P(t, z) f(z)$ for $t \in T$,
(b) $F(z)=f(z)$ on $Z$ and $F(s)=\sum_{t \in \operatorname{ims}(s)} P(s, t) F(t)$ for $s \in T \backslash Z$ (where $\operatorname{ims}(s)$ is the set of all immediate successors of $s$ in $T$ ),
(c) $F(z)=f(z)$ on $Z$ and $F(s)=\sum_{t \in L} P(s, t) F(t)$ for $s \in T \backslash Z$, whenever $L$ is a slice through $T$ below $s$.

The equivalences $(1) \Longleftrightarrow(2) \Longleftrightarrow(3)$ are now an immediate corollary of the lemma, using $(T, \sqsubset):=(\mathcal{H}, \sqsubset), P\left(h, h^{\prime}\right)=P^{\sigma}\left(h, h^{\prime}\right), f(z):=u_{p}(z)$, and $F(h):=v_{p}^{\sigma}(h)$.

### 3.2. Generalization to Information Sets

To talk about expected values of infosets, we first need to extend reach probabilities (eq. $3.1 \mid 3.9$ to infosets and define beliefs. We start by defining reach probabilities of arbitrary subsets of $\mathcal{H}$ and proving that in the specific case of infosets, reach probabilities can be factored into the contributions of individual players. Afterwards, we show players' beliefs can be defined even over infosets that are unreachable, which allows for computing reach probabilities for arbitrary paths through the infoset tree $\mathcal{I}_{p}$. This enables us to prove that infoset values' behaviour is similar to that described in Lemma 3.1 for history values (Theorem 11).

Reach Probabilities and Their Factorization. First of all, we see that the reach probability $P^{\sigma}(H)$ of a general set $H \subset \mathcal{H}$ can be computed by summing the reach probabilities of $H$ 's upper frontier [12]

$$
\begin{align*}
P^{\sigma}(H) & :=\operatorname{Pr}_{\sigma}[H \text { reached during the course of the game }]  \tag{3.14}\\
& :=\sum\left\{P^{\sigma}(h) \mid h \in H \&(\nexists g \in H, g \neq h): g \sqsubset h\right\} . \tag{3.15}
\end{align*}
$$

Assumption: We can deal with reach probabilities of infosets and public states by always taking the sum over their upper frontier. However, this would be somewhat tedious without bringing new insights. Instead, we will thus assume that all infosets and public states are thin (i.e., equal to the upper frontier). Under this assumption, 3.15 simplifies to $P^{\sigma}(I)=\sum_{h \in I} P^{\sigma}(h)$ for $I \in \mathcal{I}_{p}$, resp. $P^{\sigma}(S)=\sum_{h \in S} P^{\sigma}(h)$ for $S \in \mathcal{S}$.

As far as intuitive (yet also formal) definitions go, generalizing $p$ 's reach probability and its counterfactual counterpart to sets is straightforward, but perhaps non-obvious $\mathbb{S}^{11}$ - we define them in terms of the corresponding players attempting to maximize the chance that the set is reached:

$$
\begin{align*}
P_{-1}^{\sigma}(H) & :=\operatorname{Pr}_{\sigma}[H \text { reached } \mid \mathrm{P} 1 \text { plays to reach } H]  \tag{3.16}\\
& :=\max _{\rho_{1} \in \Sigma_{1}} P^{\rho_{1}, \sigma_{-1}}(H)  \tag{3.17}\\
P_{1}^{\sigma}(H) & :=\operatorname{Pr}_{\sigma}[H \text { reached } \mid \text { both P2 and chance play to reach } H]  \tag{3.18}\\
& :=\max _{\rho_{2} \in \Sigma_{2}} \max _{\rho_{c} \in \Sigma_{c}} P^{\sigma_{1}, \rho_{1}, \rho_{c}}(H) \tag{3.19}
\end{align*}
$$

for player 1, and similarly for player 2 . However, the resulting quantity only has desirable properties when considered over specific sets (e.g., infosets):

[^5]Lemma 3.3 (Factorization of Infoset Reach Probabilities). For any $\sigma \in \Sigma$, $p \in \mathcal{N}$, and $I \in \mathcal{I}_{p}$, we have

$$
\begin{align*}
P_{p}^{\sigma}(I) & =P_{p}^{\sigma}(h) \text { for each } h \in I,  \tag{3.20}\\
P_{-p}^{\sigma}(I) & =\sum_{h \in I} P_{-p}^{\sigma}(h), \text { and }  \tag{3.21}\\
P^{\sigma}(I) & =P_{p}^{\sigma}(I) P_{-p}^{\sigma}(I) . \tag{3.22}
\end{align*}
$$

Note that the result is very much not true in general - e.g., when $I$ is an infoset of player 1, the result doesn't hold for player 2's probabilities. Indeed, $P_{2}^{\sigma}(h)$ might be different for every $h \in I$, and the counterfactual reach probabilities $P_{-2}^{\sigma}(h), h \in I$, might sum to a number larger than 1. As with histories, we say that a set $H \subset \mathcal{H}$ is reachable (resp. unreachable) under $\sigma$ when $P^{\sigma}(H)$ is greater than zero (resp. equal to zero). An infoset of $p$ is counterfactually reachable (resp. unreachable) (by $p$ ) when $P_{-p}^{\sigma}(I)$ is greater than (resp. equal to) zero.

Beliefs. Lemma 3.3 allows us to introduce a robust notion of a belief over histories within an infoset. When $I$ is reachable, we define

$$
\begin{equation*}
P^{\sigma}(h \mid I):=\mathbf{P r}_{\sigma}[\text { current history is } h \mid \text { curr. infoset is } I]:=\frac{P^{\sigma}(h)}{P^{\sigma}(I)} \tag{3.23}
\end{equation*}
$$

When $I \in \mathcal{I}_{p}$ is the current infoset of $p, P^{\sigma}(\cdot \mid I)$ can be interpreted as $p$ 's belief about what the current state of the game is. To extend this notion even to unreachable infosets, we use an approach inspired by trembling hand equilibria, where strategies are injected with an infinitesimal amount of noise [13] (which makes $I$ reachable):

Definition 3.4 (Generalized belief over an infoset). For $I \in \mathcal{I}$ that is unreachable under $\sigma \in \Sigma$, the (generalized) belief over $I$ is $\lim _{n \rightarrow \infty} P^{\left(1-\frac{1}{n}\right) \sigma+\frac{1}{n} \text { unif }}(h \mid I)$, where unif denotes the uniformly random strategy.

On the first approximation, we can think of these beliefs as the beliefs we would obtain if all strategies that cause $I$ to be unreachable were replaced by the uniform strategy. The limit serves to give priority to the histories that only require fewer deviations from $\sigma$ to become reachable. While the choice of uniformly random noise is an ad-hoc one, the resulting concept is nevertheless compatible with the notion of belief from 3.23).

Lemma 3.5 (Equivalent definitions of the infoset belief). Let $\sigma \in \Sigma, I \in \mathcal{I}_{p}$. (1) The limit defining $P^{\sigma}(h \mid I)$ always exists. (2) For cf. reachable $I, P^{\sigma}(h \mid I)=$ $\lim _{n} P^{\left(1-\frac{1}{n}\right) \sigma+\frac{1}{n} \text { unif }}(h \mid I)=\frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}$. (3) For reachable $I$,

$$
\begin{equation*}
P^{\sigma}(h \mid I)=\lim _{n \rightarrow \infty} P^{\left(1-\frac{1}{n}\right) \sigma+\frac{1}{n} \text { unif }}(h \mid I)=\frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}=\frac{P^{\sigma}(h)}{P^{\sigma}(I)} \tag{3.24}
\end{equation*}
$$

As a corollary, it is correct to write $P^{\sigma}(h \mid I)$ without regard for $I$ 's reachability: When the infoset is reachable, the symbol can stand for $\frac{P^{\sigma}(h)}{P^{\sigma}(I)}$. When it is only reachable counterfactually, it can still stand for $\frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}$. And when unreachable even counterfactually, we have to use the limit version.

Paths in the Infoset Tree. The extension relation $\sqsubset$ (defined in Sec. 2) can be translated from $\mathcal{H}$ to $\mathcal{I}_{p}($ and $\mathcal{S})$ by saying that $J \in \mathcal{I}_{p}$ is an extension of $I \in \mathcal{I}_{p}$, written as $J \sqsupset I$, if there are some histories $h \in J, g \in I$ s.t. $h$ extends $g$. (By perfect recall, this is equivalent to saying that every element of $J$ is an extension of some element of $I$.) This turns $\mathcal{I}_{p}$ into a tree and allows us to talk about paths and slices through $\mathcal{I}_{p}$ (defined analogously to paragraph below eq. 3.9) and distinguish between terminal and non-terminal infosets (i.e., $I \subset \mathcal{Z}$ and $I \subset \mathcal{H} \backslash \mathcal{Z})$. For $I \in \mathcal{I}_{p}$,

$$
\begin{equation*}
\operatorname{ims}(I):=\left\{J \in \mathcal{I}_{p} \mid J \subset\{h a \mid h \in I, a \in \mathcal{A}(I)\}\right\} \tag{3.25}
\end{equation*}
$$

thus denotes the collection of immediate successors of $I$ in $\mathcal{I}_{p}$. For $\mathcal{P}(I)=p$,

$$
\begin{equation*}
\operatorname{ims}(I, a):=\left\{J \in \mathcal{I}_{p} \mid J \subset\{h a \mid h \in I\}\right\} \tag{3.26}
\end{equation*}
$$

denotes those immediate successors $J$ for which $a$ was chosen at $I$.
Since we already have the notion of belief $P^{\sigma}(\cdot \mid I)$, we can provide both the intuitive definition of reach probability over a path in $\mathcal{I}_{p}$ and its calculationfriendly equivalent:

$$
\begin{align*}
P^{\sigma}(I, J) & :=\mathbf{P r}_{\sigma}[J \text { reached } \mid \text { current infoset is } I],  \tag{3.27}\\
& =\sum_{g \in I} P^{\sigma}(g \mid I) \sum_{g \sqsubset h \in J} P^{\sigma}(g, h),  \tag{3.28}\\
P_{-p}^{\sigma}(I, J) & :=\operatorname{Pr}_{\sigma}[J \text { reached } \mid \text { curr. infoset is } I, p \text { plays to reach } J]  \tag{3.29}\\
& =\sum_{g \in I} P^{\sigma}(g \mid I) \sum_{g \sqsubset h \in J} P_{-p}^{\sigma}(g, h),  \tag{3.30}\\
P_{p}^{\sigma}(I, J) & :=\mathbf{P r}_{\sigma}[J \text { reached } \mid \text { curr. infoset is } I,-p \text { play to reach } J]  \tag{3.31}\\
& =\sum_{g \in I} P^{\sigma}(g \mid I) \sum_{g \sqsubset h \in J} P_{p}^{\sigma}(g, h) . \tag{3.32}
\end{align*}
$$

This guarantees that the reach probabilities of infosets can be decomposed, both in terms of the $p$-component and counterfactual component and in terms of splitting the path to any infoset into segments:

Lemma 3.6 (Properties of Infoset Reach Probabilities). For any $\sigma \in \Sigma$ and $I \sqsubset J \sqsubset K$ in $\mathcal{I}_{p}$, we have:
(1) $P^{\sigma}(I, J)=\lim _{n \rightarrow \infty} \frac{P^{\sigma^{n}}(J)}{P^{\sigma^{n}}(I)}, P_{p}^{\sigma}(I, J)=\lim _{n \rightarrow \infty} \frac{P_{p}^{\sigma^{n}}(J)}{P_{p}^{\sigma^{n}}(I)}$,
and $P_{-p}^{\sigma}(I, J)=\lim _{n \rightarrow \infty} \frac{P_{-p}^{\sigma^{n}(J)}}{P_{-p}^{\sigma^{n}}(I)}$ (where $\sigma^{n}$ denotes $\frac{n-1}{n} \sigma+\frac{1}{n}$ unif).
(2) $P^{\sigma}(I, J)=P_{p}^{\sigma}(I, J) P_{-p}^{\sigma}(I, J)$.
(3) $P^{\sigma}(I, K)=P^{\sigma}(I, J) P^{\sigma}(J, K)$.

As a result, the tree $\mathcal{I}_{p}$ satisfies the assumptions of Lemma 3.2.

Expected Utilities of Information Sets. We now have all tools required to define the expected utilities of infosets and show they behave analogously to $v_{p}^{\sigma}(h)$ :
Definition 3.7 (Infoset value). Let $I \in \mathcal{I}_{p}$ and $\sigma \in \Sigma$. The ( $p$ 's) value of $I$ under $\sigma$ is defined as

$$
\begin{equation*}
V^{\sigma}(I):=V_{p}^{\sigma}(I):=\sum_{h \in I} P^{\sigma}(h \mid I) v_{p}^{\sigma}(h) \tag{3.33}
\end{equation*}
$$

For non-terminal $I$ and $a \in \mathcal{A}_{p}(I),(p$ 's $)$ value of taking $a$ at $I$ is

$$
\begin{equation*}
Q^{\sigma}(I, a):=Q_{p}^{\sigma}(I, a):=\sum_{h \in I} P^{\sigma}(h \mid I) q_{p}^{\sigma}(h, a) \tag{3.34}
\end{equation*}
$$

By the following lemma, the action values can be derived from infoset values. As a result, we will mostly focus our analysis on the latter.

Lemma 3.8. For any $\sigma \in \Sigma, I \in \mathcal{I}_{p}$ s.t. $\mathcal{P}(I)=p$, and $a \in \mathcal{A}(I)$, we have

$$
\begin{equation*}
Q^{\sigma}(I, a)=\sum_{J \in \operatorname{ims}(I, a)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J) \tag{3.35}
\end{equation*}
$$

Since $V^{\sigma}(I)=\sum_{h \in I} P^{\sigma}(h \mid I) v_{p}^{\sigma}(h)=\sum_{h \in I} \frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)} v_{p}^{\sigma}(h)$ for counterfactually reachable $I$, the infoset values from Definition 3.7 can be viewed as a generalization of the notion of counterfactual utilities from [14] (whose definition does not work for counterfactually unreachable $I$ ).

As advertised, Theorem 1 demonstrates that infoset values can be defined in several equivalent ways. Specifically, it shows that they can be viewed as weighted sums of history-values (1) or as expected utilities calculated over the infoset tree (2-4). The latter can be viewed either as the expectation over terminal states (2), over immediate successors in the infoset tree (3, $3^{\prime}$ ), or over any slice through the infoset tree (below the given infoset) (4).
Theorem 1 (Characterization of $V^{\sigma}$ ). Suppose that terminal infosets are always singleton. Then for any $p \in \mathcal{N}$ and $\sigma \in \Sigma$, we have $V_{p}^{\sigma}$ (root $)=u_{p}(\sigma)$ and the functions $V^{\sigma}: \mathcal{I}_{p} \rightarrow \mathbb{R}$ and $Q^{\sigma}(I, a):=\sum_{\operatorname{ims}(I, a)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J)$ have the following properties:
(1) $V^{\sigma}(I)=\sum_{h \in I} P^{\sigma}(h \mid I) v_{p}^{\sigma}(h)$.
(2) $V^{\sigma}(I)=\sum_{z \in \mathcal{Z}} P^{\sigma}(I,\{z\}) u_{p}(z)$.
(3) $V_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and for non-terminal $I \in \mathcal{I}_{p}$, we have

$$
\begin{equation*}
V^{\sigma}(I)=\sum_{J \in \operatorname{ims}(I)} P^{\sigma}(I, J) V^{\sigma}(J) . \tag{3.36}
\end{equation*}
$$

(3') $V_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and for non-terminal $I$, we have

$$
V^{\sigma}(I)= \begin{cases}\sum_{a \in \mathcal{A}_{p}(I)} \sigma_{p}(I, a) Q^{\sigma}(I, a) & \text { when } p \text { acts in } I \\ \sum_{J \in \operatorname{ims}(I)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J) & \text { when } p \text { doesn't act in } I .\end{cases}
$$

(4) $V_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and for non-terminal $I \in \mathcal{I}_{p}$ and any slice $\mathcal{L}$ through $\mathcal{I}_{p}$ below $I$, we have $V^{\sigma}(I)=\sum_{J \in \mathcal{L}} P^{\sigma}(I, J) V^{\sigma}(J)$.

Moreover, each of these conditions can be used as an equivalent definition of $V_{p}^{\sigma}$ (i.e., it automatically implies all the other properties).

While the result assumes that terminal infosets are always singleton, this assumption is primarily a cosmetic one, to highlight the connection to Lemma 3.1. Indeed, if terminal infosets were of the more general form $Z \subset \mathcal{Z}$, we could equally well anchor $V^{\sigma}$ using the function $U_{p}^{\sigma}(Z):=\sum_{z \in Z} P^{\sigma}(h \mid Z) u_{p}(z){ }^{12}$

A useful corollary of Theorem 1 is that to compute $V_{p}^{\sigma}(I)$, we only need to specify $p$ 's strategy below $I$ (but neither above it nor below its siblings) and $-p$ 's strategy above and below the public state that contains $I$ (but not below the public state's siblings).

Finally, [15] recently showed that $V^{\sigma}(I)$ can also be defined as a partial derivative (or rather, more precisely, a supergradient) of the expected utility of the public state which contains $I$ with respect to the reach probability of $I$. This result thus presents another useful view of values functions, one that is orthogonal to the direction taken here. We will discuss it further in Remark 4.32 , once we have access to more related concepts.

### 3.3. Counterfactual Values and CFR

The tools developed in this section also allow us to summarize the properties of the counterfactual values used in the CFR literature [14. Since the algorithm we use in our experimental section is an extension of CFR, we use this opportunity to describe CFR's standard version.

Counterfactual Values. Counterfactual values, and the corresponding actionvalues, are defined as the non-counterfactual values weighted by the counterfactual reach probability:

$$
\begin{align*}
v_{p, \mathrm{cf}}^{\sigma}(h) & :=P_{-p}^{\sigma}(h) v_{p}^{\sigma}(h)  \tag{3.37}\\
q_{p, \mathrm{cf}}^{\sigma}(h, a) & :=v_{p, \mathrm{cf}}^{\sigma}(h a)  \tag{3.38}\\
V_{\mathrm{cf}}^{\sigma}(I) & :=V_{p, \mathrm{cf}}^{\sigma}(I) \quad:=\sum_{h \in I} v_{p, \mathrm{cf}}^{\sigma}(h),  \tag{3.39}\\
Q_{\mathrm{cf}}^{\sigma}(I, a) & :=Q_{p, \mathrm{cf}}^{\sigma}(I, a) \tag{3.40}
\end{align*}
$$

The main advantage of the counterfactual values over the "standard" ones is that they never run into problems with unreachable states - either $h$ (or $I$ ) is counterfactually reachable, and then its counterfactual value follows from the

[^6]well-defined formula (3.37), or it is cf. unreachable and its cf. value is 0 . The downside is that these values have no clear intuitive interpretation (unlike $V^{\sigma}(I)$ being the expected utility conditional on being at $I$ ).

The following result highlights the important properties of $V_{\mathrm{cf}}^{\sigma}$. First, $V_{p \mathrm{cf}}^{\sigma}($ root $)$ coincides with $u_{p}(\sigma)$ - in particular, if neither player can improve their $V_{p \mathrm{cf}}^{\sigma}($ root $), \sigma$ is an equilibrium (1). Second, $V_{\text {cf }}^{\sigma}$ can be obtained from $V^{\sigma}$ using (2). This implies (3) and (4): once we know some player's counterfactual values in leaves, backpropagating them only requires the knowledge of that player's reach probabilities. In particular, computing the cf. values of leaves is more difficult than computing their normal values but backpropagating these values is much easier.

Theorem 2 (Properties of $V_{\mathrm{cf}}^{\sigma}$ ). For any $\sigma \in \Sigma, p \in \mathcal{N}$, and non-terminal $I \in \mathcal{I}_{p}$, we have:
(1) $V_{p, \text { cf }}^{\sigma}($ root $)=u_{p}(\sigma)$.
(2) $V_{\mathrm{cf}}^{\sigma}(I)=P_{-p}^{\sigma}(I) V^{\sigma}(I)$.
(3) For any slice through $\mathcal{I}_{p}$ below $I$, we have $V_{\mathrm{cf}}^{\sigma}(I)=\sum_{J \in \mathcal{L}} P_{p}^{\sigma}(I, J) V_{\mathrm{cf}}^{\sigma}(J)$.
(4) (a) For terminal $Z \in \mathcal{I}_{p}, V_{\mathrm{cf}}^{\sigma}(Z)=\sum_{z \in Z} P_{-p}^{\sigma}(z) u_{p}(z)$.
(b) When $\mathcal{P}(I)=p$, we have

$$
V_{\mathrm{cf}}^{\sigma}(I)=\sum_{a \in \mathcal{A}(I)} \sigma_{p}(I, a) Q_{\mathrm{cf}}^{\sigma}(I, a)=\sum_{a \in \mathcal{A}(I)} \sigma_{p}(I, a) \sum_{J \in \operatorname{ims}(I, a)} V_{\mathrm{cf}}^{\sigma}(J) .
$$

(c) When $\mathcal{P}(I) \neq p$, we have $V_{\mathrm{cf}}^{\sigma}(I)=\sum_{J \in \operatorname{ims}(I)} V_{\mathrm{cf}}^{\sigma}(J)$.

Counterfactual Regret Minimization. Informally speaking, the counterfactual regret minimization algorithm (CFR) uses counterfactual values to compute an approximate equilibrium by iteratively traversing the game tree and minimizing regret at each information set [14]. Formally, the immediate counterfactual regret is the difference between counterfactual value of an infoset and the highest cf. value achievable by changing the strategy only at that infoset:

$$
\begin{equation*}
R_{\mathrm{cf}}^{\sigma}(I):=\max _{a \in \mathcal{A}(I)} R_{\mathrm{cf}}^{\sigma}(I, a):=\max _{a \in \mathcal{A}(I)} Q_{\mathrm{cf}}^{\sigma}(I, a)-V_{\mathrm{cf}}^{\sigma}(I) \tag{3.41}
\end{equation*}
$$

CFR starts with uniformly random strategy $\sigma^{1}$ and updates it using the regret matching update rule with respect to immediate cf. regrets:

$$
\sigma^{t+1}(I, a)= \begin{cases}\frac{\sum_{k=1}^{t} R_{\mathrm{cf}}^{\sigma^{k}}(I, a)}{\sum_{k=1}^{t} \sum_{a \in \mathcal{A}(I)} R_{\mathrm{cf}}^{\sigma k}(I, a)} & \text { if } \sum_{k=1}^{t} \sum_{a \in \mathcal{A}(I)} R_{\mathrm{cf}}^{\sigma^{k}}(I, a)>0  \tag{3.42}\\ \frac{1}{|A(I)|} & \text { otherwise }\end{cases}
$$

This causes the average strategy $\bar{\sigma}$ to converge to an equilibrium [14]:

$$
\begin{equation*}
\bar{\sigma}(I, a)=\sum_{t=1}^{T} \frac{P_{p}^{\sigma^{t}}(I)}{\sum_{k=1}^{T} P_{p}^{\sigma^{k}}(I)} \sigma^{t}(I, a) \tag{3.43}
\end{equation*}
$$

While CFR typically employs regret matching, it is not the only valid option. One alternative is the Exp3 algorithm [16], which always produces fully-mixed strategies. Adopting such algorithm is tempting since it makes the whole $\mathcal{H}$ is reachable and simplifies many theoretical considerations. However, this benefit comes with two limitations. First, it goes away the moment we artificially restrict the action space (as is routinely done, for example, in poker [1]). Second, practical implementations of CFR ignore the (counterfactually) unreachable parts of the game tree. Using fully-mixed strategies will therefore increase the time needed to run the algorithm.

## 4. Theory

In this section, we present the theoretical contributions made by this paper. Because of the section's length, we first give a brief summary of its contents.

We start by introducing depth-limited versions of many notions such as exploitability, best-response, and Nash equilibria (Section 4.1.1). Afterwards, we define the abstract notion of depth-limited games - these can be understood as game trees where terminal values depend not just on the leaf in which the game ends, but also on the strategy used to get to the leaf (Section 4.1.2). To utilize this notion, we need to identify depth-limited games that are useful for finding equilibria of non-depth-limited games. In Section 4.1.3, we thus provide a sufficient condition for solutions of a depth-limited game to coincide with Nash equilibria of the original game (or, more precisely, with strategies that can be extended into such Nash equilibria). Unfortunately, we also see that this criterion alone does not imply the existence of practical methods for solving the depth-limited game.

To allow for such practical methods, we define optimal value functions as those that correspond to values of optimal extensions of the trunk strategy (Section 4.2). We make a distinction between (a) reachably-optimal, (b) counterfactuallyoptimal, and (c) universally-optimal value functions, based on whether the corresponding extension of the trunk-strategy performs well (a) in situations that arise under the trunk strategy, (b) in those that could arise if one player deviated from the trunk strategy, or (c) in all situations. We look at the first two types of optimality in detail, showing how to compute the corresponding value functions and discussing which depth-limited algorithms are enabled by them for a summary, see Table 1. In particular, we prove that counterfactually-optimal value functions are both sufficient and necessary to make depth-limited CFR work well (Proposition 4.20).

In Section 4.3.1, we relax the notion of optimality to only require being optimal w.r.t. a subset of all strategies and observe that this still enables finding high-quality strategies. Moreover, we observe that [4]'s multi-valued states can be understood as a specific instance of depth-limited solving that relies on this type of optimal value functions. We argue that this connection deserves further attention.

Finally, we show how to represent value functions more compactly and prove that public states provide the minimum context that is still informative enough

| Optimality | reachable | counterfactual | universal |
| :---: | :---: | :---: | :---: |
| Described in | Section 4.2 .2 | Section 4.2 .3 | Section 4.2 .4 |
| Optimal at | $P^{\sigma^{\top}}(I)>0$ | $P_{-p}^{\sigma^{\top}}(I)>0$ | all $I$ |
| Suff. stats | joint range | range | not discussed |
| How to | value- | CFR; post-processing | not discussed |
| compute? | solving | via best-response |  |
| Algorithms | IS-MCTS, | CFR, best response, | minimax, |
| enabled | subgame value | poss. fictitious play | possibly others |

Table 1: A summary of properties of different types of optimal value functions.
to enable the computation of optimal values (Section 4.3.2). We also note that optimal values are not uniquely defined, which could prove important (and potentially inconvenient) when approximating them by neural nets (Section A). However, we do not currently see this topic as priority since we have not encountered any difficulties in practice (Section 6).

### 4.1. Depth-Limited Methods

Depth-limited methods are a standard tool in perfect-information games, and several such methods exist in imperfect information games as well 4, 17, 1, 5. Using the background given in Section 3, we now present a framework that allows for unifying these approaches.

We briefly talk about general partial strategies that are defined only on some subset of the game tree and discuss how to measure their quality. However, we quickly move to the specific case where the subset of the game tree is a so-called trunk - i.e., a connected subset that contains the tree's root. We define non-abstract value functions as mappings that take a fully-defined strategy and a node at the bottom of the trunk and return the nodes expected value under the given strategy (Section 4.1.1).

To abstract away the bottom part of the game, we define abstract value functions as mappings that output a real number for each trunk strategy and a node at the bottom of the trunk. A depth limited game is then defined analogously to a standard game, except that values of terminal nodes (at the bottom of the trunk) are determined by an abstract value function instead of a utility function. Unlike in standard imperfect-information games and depth-limited perfect-information games, the desirability of an outcome in a depth-limited imperfect-information game thus depends not just on which terminal node is generated, but also on which strategy was used to do so (Section 4.1.2p.

In any depth-limited game, we can look for optimal strategies, best responses, exploitability, etc. However, all these concepts are defined with respect to the given abstract value function - if the depth-limited game is to be useful, the value function will need to capture information about the optimal play in the part of the game that it is intended to be abstracted away. But before analyzing optimal value functions, we need to be more specific about what it means for a
depth-limited game to be useful. Initially, it might seem that all that is needed is for strategies "optimal with respect to the abstract value function" to coincide with Nash equilibria of the original game. Unfortunately, this alone isn't enough because there are value functions that satisfy this criterion while not admitting any realistic method for finding the optimal strategy. (A simple artificial example is a value function that a gives a binary penalty to any player whose strategy is exploitabl ${ }^{13}$ - the only way to find an equilibrium here is to randomly stumble upon it. However, we also later provide a realistic example.) Instead, we can talk about the usefulness in the context of specific calculations that we might wish to perform on the game: for example, does the given value function enable a depth limited version of algorithms like CFR or fictitious play? (That is, can we straightforwardly modify those algorithms to work in the depth-limited game and will their outputs converge to Nash equilibria?) Does the value function enable finding the game value of subgames? Or perhaps our interest is something other than finding Nash equilibria. For example, to allow for exploiting specific opponents, it is important that the value function enables an analogue of the best-response operator. Since the purpose of each calculation is different, the "enabling" concept is hard to capture formally. However, it can be formalized in all specific cases (Section 4.1.3)

At a high level, this section shows how to combine depth-limited games with arbitrary value functions, some of which might be of no practical use, while Section 4.2 analyzes optimal value functions, giving depth-limited algorithms guarantees analogous to their original versions.

### 4.1.1. Trunks and Partial Strategies

Before motivating and defining depth-limited games, we introduce the notion of partial strategies, trunk strategies, and their exploitability. If $p$ 's (full) strategy is function defined on $\mathcal{I}_{p}$, their partial strategy is defined on a subset of $\mathcal{I}_{p}$ :

Definition 4.1 (Partial strategy). A partial strategy of player $p$ is a mapping that assigns a probability distribution $\sigma_{p}(I) \in \Delta(\mathcal{A}(I))$ to each infoset $I$ from some subset $\mathcal{J}_{p} \subset \mathcal{I}_{p}$. Let $H \subset \mathcal{H}$ be a set closed under membership within infosets ${ }^{14} \Sigma_{p}^{H}$ denotes the set of all $p$ 's partial strategies defined on $\mathcal{J}_{p}=$ $\left\{I \in \mathcal{I}_{p} \mid I \subset H\right\}$. We set $\Sigma^{H}:=\Sigma_{1}^{H} \times \Sigma_{2}^{H}$.

When a partial strategy $\rho_{p}$ and a (possibly partial) strategy $\sigma_{p}$ satisfy (i) the domain of $\rho_{p}$ is a subset of the domain of $\sigma_{p}$ and (ii) $\sigma_{p}=\rho_{p}$ on the domain of $\rho_{p}$, we say that $\sigma_{p}$ extends $\rho_{p}$ and write $\rho_{p} \subset \sigma_{p} .^{15}$

[^7]Definition 4.2 (Trunk). A set $\mathcal{T} \subset \mathcal{H} \backslash \mathcal{Z}$ is called a trunk ${ }^{16} \mathrm{f}$ it is closed under parent nodes and membership within public states ${ }^{17}$ By $\mathcal{Z}^{\top}:=\{h a \mid h \in$ $\mathcal{T}, a \in \mathcal{A}(h), h a \notin \mathcal{T}\}$, we denote the leaves of $\mathcal{T}{ }^{18}$ The elements of $\Sigma_{p}^{\mathcal{T}}$ are called trunk strategies.

Note that the leaves of any trunk are a slice through $\mathcal{H}$ (in the sense of the definition below eq. 3.9), allowing us to apply Lemma 3.1 and Theorem 1 . In particular, the value of any history (or infoset) in $\mathcal{T}$ can be calculated by summing up the (appropriately weighted) values of its descendants in $\mathcal{Z}^{\mathcal{T}}$.

As in the full-game setting (Sec. 22), we can use exploitability to quantify the magnitude of mistakes that a partial strategy makes. Exploitability of a partial strategy is thus defined as the exploitability of the strategy's least-exploitabl ${ }^{19}$ extension to the full game $G$. To simplify the discussion, we will mostly focus on trunk strategies. However, a number of concepts we study (such as this one) could also be applied to the more general notion of partial strategies.

Definition 4.3 (Trunk exploitability). The exploitability of a trunk strategy $\sigma_{p}^{\mathcal{T}}$ is defined as

$$
\operatorname{expl}_{p}\left(\sigma_{p}^{\mathcal{T}}\right):=\min \left\{\operatorname{expl}_{p}\left(\sigma_{p}\right) \mid \sigma_{p}^{\mathcal{T}} \subset \sigma_{p} \in \Sigma_{p}\right\}
$$

We can measure the quality of trunk-strategy profiles $\sigma^{\mathcal{T}}=\left(\sigma_{1}^{\mathcal{T}}, \sigma_{2}^{\mathcal{T}}\right)$ via the average of $\operatorname{expl}_{p}\left(\sigma_{p}^{\mathcal{T}}\right)$ for $p=1,2$. We say that $\sigma^{\mathcal{T}}$ is an $\epsilon$-Nash equilibrium in the trunk, denoted as $\sigma^{\mathcal{T}} \in \epsilon$ - $\left.\mathrm{NE}(G)\right|_{\mathcal{T}}$, if it can be extended into an $\epsilon$-NE of $G$. Since Nash equilibria in $G$ are precisely the pairs of unexploitable strategies, Definition 4.3 immediately implies that equilibria in the trunk are precisely the restrictions of full-game equilibria to $\mathcal{T}$ :
Lemma 4.4. A partial strategy has $\operatorname{expl}_{p}\left(\sigma_{p}^{\mathcal{T}}\right)=0$ if and only if it can be extended into a Nash equilibrium in $G$.

Definition 4.3 also suggests that the exploitability of a partial strategy can be obtained by comparing the original value of the game $G$ with the value of the game $G\left(\sigma_{p}^{\mathcal{T}}\right)$ defined by "forcing" $p$ to play $\sigma_{p}^{\mathcal{T}}$ in $\mathcal{T}{ }^{20}$ The following formula thus gives a practical recipe for computing trunk exploitability:

[^8]Proposition 4.5 (Computing trunk exploitability). Exploitability of a trunk strategy can be computed as

$$
\operatorname{expl}_{p}\left(\sigma_{p}^{\mathcal{T}}\right)=\operatorname{gv}_{p}(G)-\operatorname{gv}_{p}\left(G\left(\sigma_{p}^{\mathcal{T}}\right)\right)
$$

Proof. For $\Sigma_{p} \ni \sigma_{p} \supset \sigma_{p}^{\mathcal{T}}$, we have $\operatorname{expl} l_{p}\left(\sigma_{p}\right)=\operatorname{gv}_{p}(G)-\min _{\sigma_{-p}} u_{p}\left(\sigma_{p}, \sigma_{-p}\right)$. Taking the minimum over $\sigma_{p}$, we get

$$
\operatorname{expl}_{p}\left(\sigma_{p}^{\mathcal{T}}\right)=\min _{\sigma_{p} \supset \sigma_{p}^{\mathcal{T}}} \operatorname{expl}_{p}\left(\sigma_{p}\right)=\operatorname{gv}_{p}(G)-\max _{\sigma_{p} \supset \sigma_{p}^{\mathcal{T}}} \min _{\sigma_{-p} \supset \sigma_{-p}^{\mathcal{T}}} u_{p}\left(\sigma_{p}, \sigma_{-p}\right)
$$

where the last term is equal to $\operatorname{gv}_{p}\left(G\left(\sigma_{p}^{\mathcal{T}}\right)\right)$.

### 4.1.2. Formal Definition of Depth-Limited Games

In the previous section on partial strategies, we have seen how to measure the quality of a trunk strategy in terms of how well would the strategy perform when extended into the full game. However, the trunk itself - without the context of the full game - is just a set of histories, not an optimization problem that can be solved. In this section, we will explain how to turn the trunk into a depth-limited game - i.e., something that can be solved - and discuss how its solutions relate to Nash equilibria in the trunk.

The central notion is that of a value function. We will be particularly interested in value functions that correspond to history-values under some strategies in the full game. However, when approximating such value functions (for example by a neural network), we might get functions which no longer correspond to any specific set of strategies, which is where the notion of abstract value functions become useful:

Definition 4.6 (Value function). Given a trunk $\mathcal{T}$, an (abstract) value function for $\mathcal{T}$ is any function $\mathbf{v}: \mathcal{Z}^{\mathcal{T}} \times \Sigma^{\mathcal{T}} \rightarrow \mathbb{R}$. We denote its values as $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)$.

When an extension $\sigma \in \Sigma$ of $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ satisfies $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)=v_{1}^{\sigma}(h)$ for every $h \in \mathcal{Z}^{\mathcal{T}}$, we say that $\mathbf{v}^{\sigma^{\mathcal{T}}}$ corresponds to $\sigma$.

In practice, value functions will typically not require the full trunk-strategy as context - for example, in perfect-information games, they require no context at all. Nor is it strictly necessary that the first component of a value-function's input is a history rather than an infoset. However, this basic case considered above is both simplest and most general. We thus present the main results in this framework and return to the topic of more compact representations of $\mathbf{v}$ in Section 4.3.2.

Definition 4.7 (Depth-limited game). For a two-player zero-sum EFG, a depthlimited game is formed by a pair $(\mathcal{T}, \mathbf{v})$, where $\mathcal{T}$ is a trunk and $\mathbf{v}$ an abstract value function.

Figure 1 shows an example of a depth-limited game. A depth-limited game is not an instance of an EFG, since $\mathbf{v}$ behaves differently than utilities do (i.e., the value of a leaf additionally depends on the strategy used to reach that leaf). However, all the traditional definitions can be re-used in $(\mathcal{T}, \mathbf{v})$ :
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-20.jpg?height=270&width=745&top_left_y=423&top_left_x=682)

Figure 1: A depth-limited version of some game. A trunk strategy is defined in the trunk $\mathcal{T}$. Trunk leaves $\mathcal{Z}^{\mathcal{T}}$ are just below the trunk, i.e., not a part of it. For a fixed trunk strategy $\sigma^{\mathcal{T}}$, the value function $\mathbf{v}^{\sigma^{\mathcal{T}}}(\cdot)$ assigns a real number to each history in $\mathcal{Z}^{\mathcal{T}}$.

Notation 4.8 (Depth-limited game concepts). Let ( $\mathcal{T}, \mathbf{v}$ ) be a depth-limited game. Expected utility in $(\mathcal{T}, \mathbf{v})$ is defined as

$$
\begin{equation*}
u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right):=\sum_{h \in \mathcal{Z} \mathcal{T}} P^{\sigma^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h), \tag{4.1}
\end{equation*}
$$

where $\mathbf{v}_{1}:=\mathbf{v}$ and $\mathbf{v}_{2}:=-\mathbf{v}$. A trunk strategy $\sigma_{p}^{\mathcal{T}} \in \Sigma_{p}^{\mathcal{T}}$ is a best response in $(\mathcal{T}, \mathbf{v})$ to $\sigma_{-p}^{\mathcal{T}} \in \Sigma_{-p}^{\mathcal{T}}$ if it maximizes the expected utility in ( $\left.\mathcal{T}, \mathbf{v}\right)$, i.e., if $\sigma_{p}^{\mathcal{T}} \in \operatorname{argmax}_{\tilde{\tilde{\sigma}}_{p}^{\mathcal{T}}} u_{p}^{\mathbf{v}}\left(\tilde{\sigma}_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)$. A trunk strategy profile is a solution $\psi^{21}$ of $(\mathcal{T}, \mathbf{v})$ if each $\sigma_{p}^{\mathcal{T}}$ is a best response in $(\mathcal{T}, \mathbf{v})$ to $\sigma_{-p}^{\mathcal{T}}$. (Similarly, we could talk about $\epsilon$-approximation solutions.) The (max-min) value of ( $\mathcal{T}, \mathbf{v})$ for player $p$ is defined as $\operatorname{gv}_{p}(\mathcal{T}, \mathbf{v}):=\max _{\sigma_{p}^{\tau}} \min _{\sigma-p} u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)$ (which stands for "game value", from reasons that will become clear later). Exploitability in ( $\mathcal{T}, \mathbf{v})$ of $\sigma_{p}^{\mathcal{T}}$ is defined as $\operatorname{expl}_{p}^{v}\left(\sigma_{p}^{\mathcal{T}}\right):=\operatorname{gv}_{p}((\mathcal{T}, \mathbf{v}))-\min _{\sigma_{-p}^{\mathcal{T}}} u_{p}^{\mathbf{v}}\left(\sigma_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)$. Finally, we use the upper index $\mathbf{v}$ to refer to ( $\mathcal{T}, \mathbf{v}$ )-variants of the various notions of value introduced in Section 3 for example, we the value of a history in $(\mathcal{T}, \mathbf{v})$ is $v_{p}^{\sigma^{\top}, \mathbf{v}}(h):=\sum_{z \in \mathcal{Z} \mathcal{T}} P^{\sigma^{\top}}(h, z) \mathbf{v}_{p}^{\sigma^{\top}}(z)$.

Until we have a reason to believe otherwise, we should not automatically assume that these "in ( $\mathcal{T}, \mathbf{v}$ )" variants of classical concepts behave as their traditional counterparts. Indeed, an arbitrary function $f: \Sigma^{\mathcal{T}} \rightarrow \mathbb{R}$ can be turned into an abstract value function via $\mathbf{v}_{f}\left(h, \sigma^{\mathcal{T}}\right):=f\left(\sigma^{\mathcal{T}}\right)$, which trivially makes the expected value of any strategy in $\left(\mathcal{T}, \mathbf{v}_{f}\right)$ equal to $u_{1}^{\mathbf{v}_{f}}\left(\sigma^{\mathcal{T}}\right)=f\left(\sigma^{\mathcal{T}}\right)$. This implies, for example, that not all $\mathbf{v}$-s will satisfy $\operatorname{gv}_{2}(\mathcal{T}, \mathbf{v})=-\operatorname{gv}_{1}(\mathcal{T}, \mathbf{v})$ (since not all functions satisfy $\max \min f=\min \max f$ ). Nevertheless, we will later show that "non-abstract" value functions (and possibly others) mostly behave as we might intuitively expect.

[^9]
### 4.1.3. Enabling Depth-Limited Algorithms

Depth-limited games and value functions are ultimately only useful to the extent that they allow us gain insights about the original game. One might intuitively expect that a value function $\mathbf{v}$ is useful if and only if the solutions of $(\mathcal{T}, \mathbf{v})$ coincide with the equilibria of the full game (Definition 4.9 below). However, we will see (in Example 4.11) that a depth-limited game can be of essentially no practical use despite satisfying this criterion. This suggests that rather than asking whether $\mathbf{v}$ is useful or not, we need a more nuanced approach: The algorithms for full EFGs naturally come with corresponding depth-limited variants. However, a depth-limited variant of an algorithm might work as intended with some value functions but not others. As a result, we propose to view a value function's usefulness in terms of the class of algorithms that it enables, i.e., that work as intended when coupled with it. We now describe these ideas in more detail and illustrate them on depth-limited minimax and CFR.

We start with the naive approach of checking for a correspondence between the solutions of ( $\mathcal{T}, \mathbf{v}$ ) and the equilibria of the full game. Recall that any trunk $\mathcal{T}$ corresponds to two types of solution concepts: First, since $\mathcal{T}$ is a part of the full game $G$, we can look at the trunk-equilibria $\mathrm{NE}(G) \mid \mathcal{T}$ given by Definition 4.3 - i.e., the strategies that can be extended to the full game in a non-exploitable manner. Second, endowing $\mathcal{T}$ with a value function $\mathbf{v}$ turns $(\mathcal{T}, \mathbf{v})$ into a depth-limited game whose solutions (described in Notation 4.8) are the strategies optimal w.r.t. $\mathbf{v}$. The latter can be computed solely based on $\mathbf{v}$ (while possibly being of no use at all), while the former are the strategies that we are after. Intuitively, one might thus expect that all we need for $(\mathcal{T}, \mathbf{v})$ to be useful is for these two notions to coincide:

Definition 4.9 (Equilibria preservation). ( $\mathcal{T}, \mathbf{v}$ ) preserves equilibria of $G$ if for every $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}, \sigma^{\mathcal{T}}$ is a solution of $(\mathcal{T}, \mathbf{v})$ if and only if $\left.\sigma^{\mathcal{T}} \in \mathrm{NE}(G)\right|_{\mathcal{T}}$.

A simple sufficient condition for equilibrium preservation is that the expected utility of any strategy $\sigma^{\mathcal{T}}$ in $(\mathcal{T}, \mathbf{v})$ coincides with value of the so-called valuesolving subgame ${ }^{22} G\left(\sigma^{\mathcal{T}}\right)$ :
Lemma 4.10. Suppose that a value function satisfies $u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$ for every $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$. Then $\mathbf{v}$ preserves Nash equilibria of $G$.

Unfortunately, equilibrium preservation is not sufficient to efficiently solve depth limited games:

Example 4.11 (The correct utility isn't all you need). Suppose we have a blackbox that takes a trunk strategy $\sigma^{\mathcal{T}}$ and returns a single scalar $\operatorname{BBox}\left(\sigma^{\mathcal{T}}\right):=$ $\operatorname{gv}_{1}\left(G\left(\sigma^{\mathcal{T}}\right)\right){ }^{23}$ This blackbox is undeniably useful, as it tells us the expected utility of the trunk strategy's "optimal extension". Additionally, setting $\mathbf{v}^{\sigma^{\mathcal{T}}}(h):=$

[^10]$\operatorname{BBox}\left(\sigma^{\mathcal{T}}\right)$ for each $h$ produces a value function which preserves the equilibria of $G$, since it (trivially) satisfies $u_{1}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\operatorname{gv}_{1}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$. However, combining this value function with algorithms that search through the game tree will be useless, since it provides no information about the individual leaves.

Now that we have the advertised negative result, we explain what is meant by "a value function $\mathbf{v}$ enabling a depth-limited version of an algorithm $A$ ".

First, before making a depth-limited modification of $A$, we need to formulate $A$ in a way that cleanly decomposes the calculation into the trunk part (to remain unchanged) and the bottom part (to be replaced by $\mathbf{v})$. For example, consider the problem of finding the value $\mathrm{gv}_{1}(G)$ of a two-player perfect-information zero-sum game $G$. The formula $\operatorname{gv}_{1}(G):=$ $\max _{\sigma_{1} \in \Sigma_{1}} \max _{\sigma_{2} \in \Sigma_{2}} u_{1}\left(\sigma_{1}, \sigma_{2}\right)$ produces the correct answer, but it is not obviously amenable to a trunk-bottom decomposition. On the other hand, the recursive calculation $\operatorname{gv}_{1}(G)=v($ root $), v(h):=$ the maximum/minimum of child values (depending on which player acts at $h$ ), $v(z):=u_{1}(z)$ for $z \in \mathcal{Z}$ is obviously decomposable - we simply replace the recursive call $v(h)$ by a call to the value function $\mathbf{v}$ when $h$ is at the bottom of the trunk. Note that for some algorithms, it is unclear (at least to the authors) how they could be usefully decomposed at all - an example is solving an EFG by a sequence-form linear program [18].

Second, a depth-limited algorithm can only be coupled with value functions that do not require more context than the algorithm internally produces. For example, any depth-limited algorithm can be coupled with $\mathbf{v}$ that doesn't require any context beyond the current history $h$ (e.g., the minimax value function $v$ from the above paragraph). On the other hand, some value functions can be much more demanding - for example, the one from Example 4.11 requires the full trunk strategy $\sigma^{\mathcal{T}}$ as context. Practical algorithms will likely fall somewhere between these two extremes. For example, in Section 4.3, we will see that depth-limited CFR (from Example 4.12 below) can be evaluated one public state $S$ at a time, requiring each player's reach probabilities of infosets in $S$.

Finally, for $\mathbf{v}$ to be useful for $A$, is not sufficient that $A$ runs when coupled with $\mathbf{v}$ - instead, $A$ coupled with v needs to achieve the purpose intended for $\mathbf{D L}-A$, which typically amounts to having guarantees analogous to $A$. For example, in perfect-information games, the minimax backpropgation algorithm can be used to obtain the minimax strategy, so our intent for DLminimax might be to find the trunk-portion of this strategy. For DL-minimax to work as we intended, it is thus sufficient to couple it with the minimax value function described above. However, other value functions might also suffice, including any $\mathbf{v}$ that values $h a$ higher than $h b$ (where $a, b \in \mathcal{A}(h)$ ) if and only if the minimax value of $h a$ is higher than the minimax value of $h b$. Note that an algorithm can have multiple purposes, which might pose different requirements on $\mathbf{v}$. For example, minimax can also be used to obtain the minimax value of the game, but DL-minimax coupled with $\mathbf{v}$ that merely preserves the action-value order would typically fail at this.

When speaking informally, we use $\mathbf{v}$ enables $\mathbf{D L}-\boldsymbol{A}$ to mean that all three
conditions above hold. While defining this notion formally for a general $A$ would be cumbersome, doing so for specific $A$ and "intent for $A$ " should be straightforward. For example, we could say that $\mathbf{v}$ enables DL-minimax iff $\mathbf{v}$ only depends on the current history (i.e., $\left.\forall h \in \mathcal{Z}^{\mathcal{T}} \forall \sigma^{\mathcal{T}}, \rho^{\mathcal{T}}: \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h)=\mathbf{v}_{p}^{\rho^{\mathcal{T}}}(h)\right)$ and the resulting DL-minimax outputs the trunk-portion of a minimax strategy in the full game.

We conclude this section by describing a depth-limited version of CFR, which will be of central interest throughout the paper.

Example 4.12 (DL-CFR). To define depth-limited CFR, we first need to formulate CFR in a way that is amenable to a trunk-bottom decomposition. Recall that standard CFR (Section 3.3) works by iteratively applying regret matching (eq. 3.42) at each infoset $I$, with respect to regrets

$$
Q_{\mathrm{cf}}^{\sigma^{t}}(I, a)-V_{\mathrm{cf}}^{\sigma^{t}}(I),
$$

where the values $V_{\mathrm{cf}}^{\sigma^{t}}(I)$ are computed by summing the terminal utilities in leaves of the full game $G$ :

$$
\begin{equation*}
V_{\mathrm{cf}}^{\sigma^{t}}(I)=\sum_{h \in I} \frac{P_{-p}^{\sigma^{t}}(h)}{P_{-p}^{\sigma^{t}}(I)} v_{p}^{\sigma^{t}}(h), \quad \text { where } v_{p}^{\sigma^{t}}(h)=\sum_{z \in \mathcal{Z}} P^{\sigma^{t}}(h, z) u_{p}(z) \tag{4.2}
\end{equation*}
$$

(and similarly for $Q_{\mathrm{cf}}^{\sigma^{t}}(I, a)$ ). To make 4.2 amenable to decomposition, we express it in terms of the values of the leaves of the trunk:

$$
\begin{align*}
v_{p}^{\sigma^{t}}(h) & =\sum_{z \in \mathcal{Z}} P^{\sigma^{t}}(h, z) u_{p}(z)=\sum_{z^{\prime} \in \mathcal{Z}^{\mathcal{T}}} \sum_{z \in \mathcal{Z}} P^{\sigma^{t}}\left(h, z^{\prime}\right) P^{\sigma^{t}}\left(z^{\prime}, z\right) u_{p}(z)  \tag{4.3}\\
& =\sum_{z^{\prime} \in \mathcal{Z}^{\mathcal{T}}} P^{\sigma^{t}}\left(h, z^{\prime}\right) v_{p}^{\sigma^{t}}\left(z^{\prime}\right) \tag{4.4}
\end{align*}
$$

Plugging $\mathbf{v}$ into (4.4), the computation only needs to reach the trunk leaves:

$$
V_{\mathrm{cf}}^{\sigma^{t, \mathcal{T}}, \mathbf{v}}(I):=\sum_{h \in I} \frac{P_{-p}^{\sigma^{t, \mathcal{T}}}(h)}{P_{-p}^{\sigma^{t, \mathcal{T}}}(I)} v_{p}^{\sigma^{t, \mathcal{T}}, \mathbf{v}}(h), \text { where } v_{p}^{\sigma^{t, \mathcal{T}}, \mathbf{v}}(h) \sum_{z^{\prime} \in \mathcal{Z}^{\mathcal{T}}} P^{\sigma^{t}}\left(h, z^{\prime}\right) \mathbf{v}_{p}^{\sigma^{t}}\left(z^{\prime}\right)
$$

(and similarly for $Q_{\mathrm{cf}}^{\sigma^{t}, \mathbf{v}}(I, a)$ ). Finally, we define DL-CFR as an algorithm that works identically to standard CFR (Section 3.3), except it only computes strategy in the trunk and performs regret matching updates with respect to $Q_{\mathrm{cf}}^{\sigma^{t}, \mathbf{v}}(I, a)-V_{\mathrm{cf}}^{\sigma^{t, \mathcal{T}}, \mathbf{v}}(I)$ instead of $Q_{\mathrm{cf}}^{\sigma^{t}}(I, a)-V_{\mathrm{cf}}^{\sigma^{t}}(I)$. If we wanted to formalize "v enabling DL-CFR", we could interpret it, for example, as "the exploitability of DL-CFR's average (trunk) strategy goes to 0 as the number of iterations increases".

In the next section, we investigate a class of value functions which are in some sense optimal, allowing DL-CFR and other algorithms to work as intended.

### 4.2. Optimal Value Functions

In the previous section, we argued that not all abstract value functions convey a sufficient amount of information about the game. So, which value functions are useful? We start by discussing the obvious and elegant answer (previously adopted by [5] and [15, Thm. 1-3]): look at value functions that correspond to Nash equilibria in subgames (Section 4.2.2). We initially believed that this is where the paper's story ends. Unfortunately, this is not the case. If our goal was only to compute the game-values of subgames (as in [5]), these value functions would be sufficient (Proposition 4.17). However, Example 4.18 shows that these value functions can fail to be useful for practical algorithms such as CFR or fictitious play.

Unfortunately, this is not the case, as this answer turns out to be insufficient for many purposes. However, such value functions can fail to be useful for practical algorithms such as CFR or fictitious play (Example 4.18). Fortunately, value functions can be useful for these purposes if they satisfy additional properties (Section 4.2.3). ${ }^{24}$

To give a better answer, we introduce the notion of optimality with respect to a collection of infosets (Section 4.2.1). Informally speaking, an extension of a trunk strategy is optimal with respect to an infoset if the infoset's owner cannot improve their expected utility at that infoset by changing their strategy in the bottom of the game. A strategy's extension is better if it is optimal with respect to a bigger collection of infosets. And a value function is better if it corresponds to better extensions. In this terminology, the initial wrong answer proposed to use value functions that are optimal only with respect to infosets that are reachable (i.e., have non-zero probability of being encountered under the given trunk strategy).

The most practical alternative answer we have identified is to use a stronger notion of counterfactual optimality, where players maximize infoset values for all infosets that their opponent's strategy allows them to enter. (That is, those that are reachable by the given player in the counterfactual scenario where the player decides to do so; see Section 4.2.3) We show that obtaining value functions of this type is not difficult - we can either use CFR to find the Nash equilibria of subgames or apply a simple post-processing step to Nash equilibria found by an arbitrary method. (The former is done in [15].) Moreover, they enable a depth-limited variant of CFR, best response, and likely also fictitious play. We also describe a simple example where depth-limited CFR and fictitious play both fail when combined with a value function that is reachably optimal but not counterfactually optimal.

Finally, we also briefly discuss an even stronger notion of universally-optimal value functions which maximize infoset values in all infosets, including those

[^11]reachable only when both the player and their opponent deviate from their strategy.

### 4.2.1. Optimality Criterion

In this section, we describe the optimal value functions promised above. The starting point is the notion of optimality with respect to a single infoset, which states that the infoset's owner cannot improve the infoset's expected utility by switching to a different strategy:

$$
\begin{equation*}
V_{p}^{\sigma}(I)=\max _{\rho_{p} \in \Sigma_{p}} V_{p}^{\rho_{p}, \sigma_{-p}}(I)=: V_{p}^{*, \sigma_{-p}}(I), \text { where } I \in \mathcal{I}_{p} \tag{4.5}
\end{equation*}
$$

Before proceeding with the formal definitions, Remark 4.13 provides further intuition regarding the usage of this criterion.

Remark 4.13 (Understanding infoset optimality). Firstly, we see that the maximum in 4.5 is taken over all strategies of $p$. However, recall that as far as $p$ 's strategy is concerned, $V_{p}^{\rho_{p}, \sigma_{-p}}(I)$ only depends on the restriction of $\rho_{p}$ to the part $\operatorname{Desc}(I):=\left\{J \in \mathcal{I}_{p} \mid J \sqsupset I\right\}$ of $p$ 's infoset tree that lies below $I$ (including $I$; this follows from Theorem 11). As a result, it doesn't matter whether the maximum is taken over the whole $\Sigma_{p}$, over $\Sigma_{p}^{\operatorname{Desc}(I)}$, or over any set inbetween.

Second, which infosets should be required to satisfy 4.5)? Since our goal is to design value functions for depth-limited solving, we can start by restricting the condition to the infosets that lie in $\mathcal{Z}^{\mathcal{T}}$. Among all possible extensions $\sigma$ of $\sigma^{\mathcal{T}}$, it is - all else being equal - better if $\sigma$ is optimal with respect to as many $I$ 's as possible. Indeed, we can view $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)$ as providing information about $h$ - if $h$ is unreachable under $\sigma^{\mathcal{T}}$, we might nevertheless find uses for any pieces of "off-policy data" about it. Another useful metaphor is to view game-solving as following a gradient of strategy quality - if the extension $\sigma$ is optimal w.r.t. more infosets, the gradient will convey more information and we can learn faster. However, the assumption of "all else being equal" does not hold in practice. Indeed, the goal will typically be to obtain the easiest-to-find extension $\sigma \supset \sigma^{\mathcal{T}}$ which is informative enough for our method of choice.

The following definition gives a general "blueprint" for various versions of optimality and introduces the three types that will be the most relevant for us. Informally speaking, it simply states that reachably (resp. counterfactually, resp. universally) optimal extensions of a trunk strategy are those for which $V_{p}^{\sigma}(I)=V_{p}^{*, \sigma_{-p}}(I)$ holds in all reachable (resp. counterfactually-reachable, resp. all) infosets ${ }^{25}$ at the bottom of the trunk.

Definition 4.14 (Optimal extensions). Let $\mathcal{J} \subset\left\{I \in \mathcal{I} \mid I \subset \mathcal{Z}^{\mathcal{T}}\right\}$ be collection of leaf-infosets for a trunk $\mathcal{T}$ and $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$. An extension $\sigma$ of $\sigma^{\mathcal{T}}$ is said to be optimal on $\mathcal{J}$, denoted $\sigma \in \operatorname{OE}\left(\sigma^{\mathcal{T}}, \mathcal{J}\right)$, if $V^{\sigma}(I)=V^{*, \sigma_{-p}}(I)$ holds for every $I \in \mathcal{J}$.

[^12]The extension is said to be reachably-optimal if it is optimal on $\mathcal{J}=\{I \subset$ $\mathcal{Z}^{\mathcal{T}} \mid I \in \mathcal{I}_{p}$ reachable $\}$. Analogously, we define counterfactually-optimal and universally-optimal extensions as those corresponding to $\mathcal{J}=\left\{I \subset \mathcal{Z}^{\mathcal{T}} \mid\right.$ $I \in \mathcal{I}_{p}$ counterfactually reachable $\}$, resp. $\mathcal{J}=$ all infosets in $\mathcal{Z}^{\mathcal{T}}$.

As we indicated earlier, any notion of optimality of extensions can be derived from the corresponding notion of optimality of values. Informally speaking, reachably-optimal value functions are those that correspond to reachably-optimal extensions on reachable infosets (but can take arbitrary values elsewhere), counterfactually-optimal value functions correspond to cf.-optimal extensions on cf.-reachable infosets (but have arbitrary values elsewhere), and universally-optimal value functions are those that corresponds to universally-optimal extensions (on the whole $\mathcal{Z}^{\mathcal{T}}$ ):

Definition 4.15 (Optimal value functions). Let $\mathcal{T}$ be a trunk and $\overrightarrow{\mathcal{J}}=\left(\mathcal{J}^{\sigma^{\mathcal{T}}}\right)_{\Sigma \mathcal{T}}$, $\mathcal{J}^{\sigma^{\mathcal{T}}} \subset\left\{I \in \mathcal{I} \mid I \subset \mathcal{Z}^{\mathcal{T}}\right\}$. A value function $\mathbf{v}: \mathcal{Z}^{\mathcal{T}} \times \Sigma^{\mathcal{T}} \rightarrow \mathbb{R}$ is $\overrightarrow{\mathcal{J}}$-optimal if it "corresponds to extensions that are optimal on $\bigcup \mathcal{J}^{\sigma^{\mathcal{T}}}$ " in the following sense:

$$
\left(\forall \sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}\right)\left(\exists \sigma \in \operatorname{OE}\left(\sigma^{\mathcal{T}}, \mathcal{J}^{\sigma^{\mathcal{T}}}\right)\right)\left(\forall h \in \bigcup \mathcal{J}^{\sigma^{\mathcal{T}}}\right): \mathbf{v}^{\sigma^{\mathcal{T}}}(h)=v_{1}^{\sigma}(h)
$$

In Sections 4.2.2 4.2.4, we investigate these different variants of optimality in the order of increasing strictness. We explain how to compute the (first two types of) value functions, lists some of the depth-limited algorithms they enable, and explain some of their shortcomings.

### 4.2.2. Reachably-Optimal Value Functions

Reachably-optimal value functions are the easiest to obtain, but also only "enable" fewer depth-limited computations than the other two types of optimal value functions. The following result shows that reachably-optimal values can be computed by fixing the trunk strategy and solving the remainder of the game:
Proposition 4.16 (Computing reachably-optimal values). Suppose that for each $\sigma^{\mathcal{T}}$, there is some $\sigma \in \operatorname{NE}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$ s.t. $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)=v_{1}^{\sigma}(h)$ holds for all $h \in \mathcal{Z}^{\mathcal{T}}$. Then $\mathbf{v}$ is a reachably optimal value function.

The proof goes by showing that if an extension of $\sigma^{\mathcal{T}}$ isn't reachably optimal, one of the players can gain by deviating from it. Note that in the other direction, it can also be shown that all reachably-optimal value functions are of this form, except for the fact that they can be defined arbitrarily on the unreachable parts of $\mathcal{Z}^{\mathcal{T}}$.

Proposition 4.17 shows that as long as we are only after computing the utility of $\sigma^{\mathcal{T}}$ 's best extension (i.e., the subgame value), reachably-optimal value functions are sufficient.

Proposition 4.17 (Enabling utility calculation). Any reachably-optimal value function satisfies

$$
\begin{equation*}
u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\sum_{h \in \mathcal{Z} \mathcal{T}} P^{\sigma^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h)=\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right) . \tag{4.6}
\end{equation*}
$$

This can be shown by combining Theorem 1 with the fact that expected utility of a strategy only depends on the values of reachable histories. In particular, this proposition demonstrates that for any trunk-strategy, reachably-optimal value functions enable us to compute the expected utility (at the root) obtained if both players play "optimally" in the remainder of the game. (Because no matter which version of optimality we choose, this number will be equal to $\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$.) Together with Lemma 4.10, Proposition 4.17 immediately yields the following corollary:

Theorem 3 (Equilibrium preservation). Reachably-optimal value functions preserve Nash equilibria.

By Theorem 3, it is also true that the equilibria of the corresponding depthlimited game $(\mathcal{T}, \mathbf{v})$ are precisely the equilibria of $G$ restricted to $\mathcal{T}$. Unfortunately, this is not the same as offering a practically viable method for finding these equilibria. Indeed, as shown by Example 4.18, reachably-optimal value functions might fail to enable depth-limited variants of more efficient solution methods such as CFR (and fictitious play, as argued later):

Example 4.18 (Reachably-optimal value functions do not enable DL-CFR). Let $G$ be a two-round simultaneous-move (but otherwise perfect-information) game that works as follows: The first round looks like matching pennies ${ }^{26}$, except that the utility for heads-heads is $\pm(1+\epsilon)$ instead of 1 (this is to break the symmetry while running CFR). Only player 1 acts in the second round and has a choice between 'doing nothing' ( $d n$; no change to utilities) or 'going mad' ( gm ; transfer 1000 utility to the opponent). The only Nash equilibrium of $G$ is to 'do nothing' in the second round and play something very close to the uniform strategy in the first round.

Consider the trunk $\mathcal{T}:=$ the first round, for which each information set just below $\mathcal{T}$ corresponds to an element of $\mathcal{Z}^{\mathcal{T}}=\{(H, H),(H, T),(T, H),(T, T)\}$. For any trunk strategy $\sigma^{\mathcal{T}}$, the only extension that satisfies $V_{p}^{\sigma}(I)=V_{p}^{*, \sigma_{-p}}(I)$ for every $I$ (independently of its reachability) is the $\sigma$ that 'does nothing' in the second round of the game, independently of what has happened in the trunk. However, if we only require optimality in reachable states, $\sigma$ is suddenly allowed to 'go mad' in states that are only reachable counterfactually. In particular, the following value function is reachably optimal:

$$
\begin{aligned}
& \mathbf{v}^{\sigma^{\tau}}(h):=u_{1}(h, g m) \text { whenever } P_{1}^{\sigma^{\tau}}(h)=0, \\
& \mathbf{v}^{\sigma^{\tau}}(h):=u_{1}(h, d n) \text { otherwise. }
\end{aligned}
$$

In other words, $\mathbf{v}$ corresponds to player 1 always believing that they would go mad if they deviated from their current trunk strategy.

We claim that running DL-CFR on ( $\mathcal{T}, \mathbf{v})$ produces a strategy where player 1 chooses $H$ with (near) certainty - i.e., one that is very far from being a Nash

[^13]equilibrium. To see this, note that DL-CFR starts out with a uniform trunk strategy, for which $\mathbf{v}$ predicts the 'do nothing' action everywhere in the bottom. The first CFR update thus works as if we were running CFR on the one-round game. Since the utility for $(H, H)$ is slightly above 1 , the maximizing player will play pure $H$ in the second iteration, while the minimizing player 2 will play pure $T$. Since $(H, T)$ is an optimal outcome for player 2 , they will only deviate from playing $T$ if player 1 deviates from $H$. However, given how $\mathbf{v}$ is defined, this will never happen. Indeed, while player 1 sees that the $(H, T)$ outcome yields -1 utility, they believe that playing $(T, T)$ would result in $u_{1}(T, T, g m)=-999$ utility. The average strategy of DL-CFR in $(\mathcal{T}, \mathbf{v})$ will thus converge to the highly-exploitable strategy $(H, T)$.

However, despite not being sufficiently informative for CFR, reachablyoptimal value functions will nevertheless be suitable for other algorithms. For example, they are a natural candidate for obtaining a depth-limited variant of information-set MCTS (IS-MCTS) [19. The standard IS-MCTS algorithm works by running some MCTS algorithm on $\mathcal{H}$ while sharing statistics between all histories within every infoset. (The algorithm is far from state of the art and lacks theoretical guarantees, so it is best used as a cheap baseline solution.) The corresponding depth-limited IS-MCTS algorithm can thus be obtained by stopping the MCTS recursion in the trunk leaves and using the values $\mathbf{v}^{(\cdot)}(h)$ corresponding to the current strategy ${ }^{27}$.

We could also consider evolutionary methods for playing imperfect information games: for example, we could have a population of agents play each other and keep the successful ones, modifying them slightly before the next iteration. All that is needed to execute a depth-limited version of this method is to obtain the correct expected utilities at the root. By Proposition 4.17, reachably-optimal value functions are sufficient for the task.

### 4.2.3. Counterfactually-Optimal Value Functions

We now turn to the stronger notion of counterfactual optimality. Informally speaking, counterfactually-optimal value functions give a "correct answer" to the question "so, this information set that I didn't play into, what would have happened if I did play there?". (Which was not the case for reachably-optimal value functions.) As we will see, this ensures that using such functions together with depth-limited CFR produces provably correct strategies. We also show how cf.-optimal value functions enable the depth-limited computation of best-response and explain how this could be useful for enabling depth-limited fictitious play.

As observed in [20, counterfactual optimality can be achieved by postprocessing the Nash equilibria of value-solving subgames (described below Def. 4.9). This can be done via "counterfactual best-response", i.e., a recursively calculated pure strategy $c b r_{p} \in \Sigma_{p}$ that takes some action from

[^14]$\operatorname{argmax}_{a \in \mathcal{A}(I)} Q^{c b r_{p}, \sigma_{-p}}(I, a)$ in each $I \in \mathcal{I}_{p}$ that is counterfactually reachable by $p$.

Proposition 4.19 (Computing counterfactually-optimal values). Suppose that for each $\sigma^{\mathcal{T}}$, a value function $\mathbf{v}$ is of the form $\mathbf{v}^{\sigma^{\top}}(h)=v_{1}^{\mu}(h)$, where $\mu \in \Sigma$ is obtained by

- starting with some $\sigma \in \operatorname{NE}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$,
- for both $p=1,2$, going through all $I \subset \mathcal{Z}^{\mathcal{T}}, I \in \mathcal{I}_{p}$, that are counterfactually reachable by $p$ but not reachable,
- and replacing $\sigma_{p}$ by $\operatorname{cbr} r_{p}\left(\sigma_{-p}\right)$ on such infosets and their descendants.

Then $\mathbf{v}$ is counterfactually optimal.
A particular consequence (of the proof) is that every reachably-optimal value function can be made counterfactually optimal by "fixing" the values in infosets that are only reachable counterfactually. In practice, we can also obtain counterfactually-optimal value functions by running CFR on the bottom of the game ${ }^{28}$ This produces a NE of $G\left(\sigma^{\mathcal{T}}\right)$ in which the values in all counterfactually reachable infosets are already maximal [14], so the post-processing step will be unnecessary.

Depth-Limited CFR. In Example 4.12, we have described the depth-limited variant of CFR. With counterfactually-optimal value functions, we are now ready to prove that DL-CFR converges to a Nash equilibrium of the game.

Informally, CFR works by finding a strategy in which both players maximize their counterfactual values $P_{-p}^{\sigma}(I) V_{p}^{\sigma}(I)$, which is equivalent to maximizing the values $V_{p}^{\sigma}(I)$ in counteractually reachable infosets. In comparison, DL-CFR works by (i) explicitly minimizing the (cf.) regret in the trunk and (ii) implicitly minimizing the (cf.) regret in the remainder of the game (by using a value function that corresponds to a zero-regret strategy). One way to prove the correctness of DL-CFR is thus to formalize this idea and mimick the standard CFR proof. However, since the existing CFR-D algorithm ${ }^{29}$ 17] also relies on a similar idea, we can show that our DL-CFR approach fits the CFR-D formalism, which guarantees that the existing CFR-D results apply to DL-CFR as well:

Proposition 4.20 (Enabling DL-CFR). Let ( $\mathcal{T}$, v) be a depth-limited game corresponding to a counterfactually-optimal value function. Then:
(1) $D L-C F R$ can be viewed an instance of $C F R-D$ and inherits its guarantees.

[^15](2) In particular, the strategy $\bar{\sigma}^{\mathcal{T}, t}$ produced aftert iterations of $D L-C F R$ satisfies $\operatorname{expl}\left(\bar{\sigma}^{\mathcal{T}}, t\right) \xrightarrow{t \rightarrow \infty} 0$.
In general, the same is not true for reachably-optimal value functions.
Depth-Limited Best-Response and Fictitious Play. We now show that counterfactually optimal value functions also enable an efficient computation of best response. In Notation 4.8 we have defined a best response to $\sigma_{-p}^{\mathcal{T}}$ in $(\mathcal{T}, \mathbf{v})$ as the strategy $\sigma_{p}^{\mathcal{T}} \in \Sigma_{p}^{\mathcal{T}}$ that maximizes the expression
$$
u_{p}^{\mathbf{v}}\left(\mu_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)=\sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\mu_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\mu_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}}(h)
$$

However, this type of best response is difficult to calculate, because the optimal course of action in the bottom part of the game depends on the trunk strategy, which in turn depends on the bottom part ${ }^{30}$ We could instead use a more efficient but less rigorous approach which pretends that modifying the trunk strategy has no influence on the output of the value function:
Definition 4.21 (Naive best-response). A trunk strategy $\sigma_{p}^{\mathcal{T}}$ is a naive bestresponse to a trunk strategy profile $\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)$ in $(\mathcal{T}, \mathbf{v})$ if it satisfies

$$
\sigma_{p}^{\mathcal{T}} \in \underset{\mu_{p}^{\mathcal{T}}}{\operatorname{argmax}} \sum_{h \in \mathcal{Z}^{\mathcal{T}}} P^{\mu_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}}(h)
$$

Because of its dependency on $\rho_{p}^{\mathcal{T}}$, a naive best-response might not be a true best response. However, for counterfactually optimal values, this concept nevertheless preserves a crucial property of standard best responses:

Proposition 4.22 (NE as mutual naive best-responses). If $\mathbf{v}$ is a a counterfactually optimal value function, then any $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ where each $\sigma_{p}^{\mathcal{T}}$ is a naive best-response to $\sigma^{\mathcal{T}}$ in $(\mathcal{T}, \mathbf{v})$ can be extended into a Nash equilibrium in $G$.

As with DL-CFR, the result doesn't hold if $\mathbf{v}$ is only reachably optimal (as witnessed using the counterexample constructed in Example 4.18.

Another algorithm that might be particularly suited for depth-limited modification is fictitious play (FP) [21, which can be formulated as a process that keeps a growing list of strategies, where each newly added strategy is a best response to the average of the strategies identified thus far (and that returns this average when terminated). Replacing best response by naive best response will therefore produce depth-limited fictitious play, which should work analogously to standard FP. Unfortunately, there are games in which fictitious play fails to find a strong strategy even when given infinite time. From this reason, we defer the analysis of its depth-limited variants to a future work. However, we believe that to the extent that FP does have some performance guarantees, DL-FP should be combined with counterfactually optimal value functions to preserve them.

[^16]
### 4.2.4. Universally-Optimal Value Functions

Universal optimality is the strongest of the tree variants of optimality introduced in Definition 4.15- universally-optimal value functions correspond to strategies whose values are optimal for all infosets in $\mathcal{Z}^{\mathcal{T}}$, including those that are completely unreachable under the given trunk strategy. While this notion often unnecessarily strong, it allows us to make an important connection between our analysis and value functions in perfect information games. Indeed, recall that perfect information games have a canonical a value function - the one that assigns to each state its minimax value. Since reachably-optimal and counterfactually-optimal values can be arbitrary in (counterfactually) unreachable parts of $\mathcal{Z}^{\mathcal{T}}$, neither of them coincides with minimax values when applied to perfect information games. In contrast, universally-optimal value functions do successfully generalize minimax values since they are optimal everywhere. However, since counterfactual optimality is sufficient for the purposes of this paper, we will not discuss universal optimality in more detail.

### 4.3. Efficient Computation of Optimal Value Functions

Our analysis so far has focused on theoretical properties of optimal value functions, leaving aside the surrounding computational considerations. In this section, we show how to relax and generalize some of the formal definitions to make dealing with optimal value functions more practical. We also analyze properties of these functions that are relevant to their approximation.

One formal limitation of our earlier definitions is that optimal value functions require the player's infoset values to be at least as high as values under all other strategies (Section 4.3.1). This is both difficult to verify and unnecessary - in practice, we might do equally well by relaxing the assumption to only include some representative subset of strategies (e.g., pure undominated strategies). By pruning the considered strategy set even further, we can make trade-offs between computational costs and theoretical guarantees. We show that this approach can be implemented in a way that reduces to the existing "multivalued states" method described in 4. This exposes a previously unrecognized connection between two coexisting depth-limited approaches to solving imperfect information games, 1] and [4], and opens up the potential for their cross-fertilization.

A second formal limitation is that value functions take the whole trunk strategy as an input and return a separate value for each history at the bottom of the trunk (Section 4.3.2). To decrease the size of the input and localize the computation of values, we show that public belief states provide a sufficient context for computing optimal values. Moreover, we show that information about anything less than a public state might fail at this task. We also argue that for practical purposes, the output should be aggregated over infosets, as is often done in the recent literature [1, 15].

Thirdly, we remark that in practice, it will be desirable to use value functions that are only $\epsilon$-optimal, typically in conjunction with neural networks. For this purpose, it would be desirable to have a unique approximation target. However, it turns out that optimal value functions are not uniquely defined. We thus
argue that the next best thing to hope for is that the set of such functions is convex. We show that this is true for reachably-optimal value functions (this is both straightforward and already known). As an open problem, we ask whether the same is true for counterfactual and universal optimality - we provide some evidence both for and against this claim. Since these details might be of interest to fewer readers, they are presented in A

### 4.3.1. Alternative Computation of Optimal Value Functions and the Relation to Multi-Valued States

In this section, we show an alternative method of computing value functions, connecting the value-function approach to depth-limited search with the multivalued states introduced in (4).

To motivate the approach, let us inspect the condition which appears in the definition of optimality (Definition 4.14). To count as a (reachably) optimal extension of a trunk strategy $\sigma^{\mathcal{T}}$, a strategy $\sigma$ needs to satisfy the following (for every reachable trunk-infoset of $p, p=1,2)$ :

$$
\left(\forall \rho_{p} \in \Sigma_{p}\right): V^{\sigma_{p}, \sigma_{-p}}(I) \geq V^{\rho_{p}, \sigma_{-p}}(I)
$$

However, taken literally, this condition requires comparing $\sigma_{p}$ to all other strategies $\rho_{p} \in \Sigma_{p}$, and there are uncountably many of those. Instead, we can consider a weaker - but often sufficient - requirement that $\sigma$ performs at least as well as any strategy from some portfolio of strategies $\mathbb{P} \subset \Sigma^{\mathcal{H} \backslash \mathcal{T}}$. The following definition captures this idea formally:

Definition 4.23 (Optimality w.r.t. a portfolio). Let $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ and $\mathbb{P}=\mathbb{P}_{1} \times \mathbb{P}_{2} \subset$ $\Sigma^{\mathcal{H} \backslash \mathcal{T}}$. An extension $\sigma$ of $\sigma^{\mathcal{T}}$ is said to be reachably optimal with respect to $\mathbb{P}$ if for every reachable trunk-leaf infoset $I \in \mathcal{I}_{p}$,
(i) $\sigma$ below $I$ is realizable by $\mathbb{P}$ (i.e., the restriction of $\sigma_{p}$ to $\left\{J \in \mathcal{I}_{p} \mid J \sqsupset I\right\}$ is a convex combination of elements of $\mathbb{P}_{p}$ ) and
(ii) for every $\rho_{p}^{\downarrow} \in \mathbb{P}_{p}, V^{\sigma_{p}, \sigma_{-p}}(I) \geq V^{\rho_{p}^{\downarrow}, \sigma_{-p}}(I) 3^{31}$

A value function $\mathbf{v}$ is reachably optimal w.r.t. $\mathbb{P}$ if for each $\sigma^{\mathcal{T}}, \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(\cdot)=v_{p}^{\sigma}(\cdot)$ for some extension $\sigma$ of $\sigma^{\mathcal{T}}$ that is reachably optimal w.r.t. $\mathbb{P}$. Counterfactual and universal optimality w.r.t. $\mathbb{P}$ is defined analogously.

In the extreme case where $\mathbb{P}$ contains all strategies, the notion trivially coincides with reachable optimality. However, we might get a high-quality value function even with a much smaller portfolio: To start with, we can reduce $\mathbb{P}$ without hurting the quality of $\mathbf{v}$ at all. Indeed, Proposition 4.24 provides a theoretical lower bound on the degree of such reduction, and a much more extensive (yet still lossless) reduction might be possible in practice:

[^17]Proposition 4.24 (Inspired by [4]). Suppose that $\mathbf{v}$ is reachabll ${ }_{32}$ optimal w.r.t. a portfolio that contains (the trunk restrictions of) all pure undominated strategies. Then $\mathbf{v}$ preserves the equilibrib of $G$.

Further reduction in the portfolio's size can come at some cost to the quality of $\mathbf{v}$. However, we might be able to reduce $\mathbb{P}$ down to a handful of strategies while only harming v a little. For example, 4] constructed a competitive agent for heads-up no-limit Texas hold'em poker using less than ten strategies ${ }^{34}$.

A practical method of computing optimal values can be devised by considering what we might call a "partial normal-form representation" of the game. The underlying idea is that instead of taking all actions sequentially, the players can only play sequentially until they reach the bottom of the trunk, at which point they simultaneously ${ }^{35}$ announce the strategy they would use in the bottom part of the game. The game then immediately terminates, giving each player the expected utility corresponding to combination of the selected strategies. When the choice of the bottom strategy is restricted to the portfolio, this corresponds to the following game:

Definition 4.25 (Inspired by [4]). For a trunk $\mathcal{T}$ and portfolio $\mathbb{P}, G(\mathcal{T}, \mathbb{P})$ is a game which proceeds identically to $G$, up to the point where it reaches some trunk-leaf $z \in \mathcal{Z}^{\mathcal{T}}$. The information sets on $\mathcal{Z}^{\mathcal{T}}$ are the same as in $G$. In $z \in \mathcal{Z}^{\mathcal{T}}$, the players simultaneously select some $\sigma_{p}^{\downarrow} \in \mathbb{P}_{p}$. The resulting history $z \sigma_{1}^{\downarrow} \sigma_{2}^{\downarrow}$ is terminal in $G(\mathcal{T}, \mathbb{P})$ and yields utility $u_{p}\left(z \sigma_{1}^{\downarrow} \sigma_{2}^{\downarrow}\right):=v_{p}^{\sigma_{1}^{\downarrow}, \sigma_{2}^{\downarrow}}(z)$.

As advertised, the equilibria of the game $G(\mathcal{T}, \mathbb{P})$ coincide with the solutions of the depth-limited game $(\mathcal{T}, \mathbf{v})$ whose value function corresponds to $\mathbb{P}$ :

Theorem 4. For any $\mathcal{T}, \mathbb{P}$, and $\mathbf{v}$ that is reachably optimal w.r.t. $\mathbb{P}$, a trunk strategy $\sigma^{\mathcal{T}}$ is solution of $(\mathcal{T}, \mathbf{v})$ if and only it is in $\left.\operatorname{NE}(G(\mathcal{T}, \mathbb{P}))\right|_{\mathcal{T}}$.

Finally, Theorem 4 has an important implication for the existing research on depth-limited solving (DLS): In the past, one line of DLS research was based on value functions (primarily building on [1). Around the same time, 4] suggested a second approach called multi-valued states. In our terminology, this approach consists of using the game $G(\mathcal{T}, \mathbb{P})$ with $\mathbb{P}_{1}:=$ all pure strategies (resp. a small set of heuristic strategies in practice) and $\mathbb{P}_{2}$ which consists of a single "blueprint

[^18]strategy" (typically obtained by solving an abstraction of $G$ ). Initially, these approaches might seem incompatible or unrelated - they certainly did to the authors of this text. However, with Theorem 4, we see that they constitute two different methods of approximating the same value function. Moreover, now that the connection has been made explicit, we can make progress by combining the ideas behind the two approaches. For example, when approximating a value function by a neural network, we can use the games $G(\mathcal{T}, \mathbb{P})$ to generate the training data.

### 4.3.2. Compact Representation of Optimal Value Functions

In this section, we discuss how to represent the input and output of optimal value functions more compactly. These improvements are particularly useful in practical applications, where value functions are often approximated by neural networks.

Input. The input can be compressed in two steps: first, we can "flatten" the trunk strategy by only considering the corresponding reach probabilities. Second, we can "narrow it down" by only considering reach probabilities that are relevant for computing the value of the specific history. The first part is formally captured by the following definition (inspired by [5]) and its simple corollary:

Definition 4.26 (Sufficient statistic). A quantity $s: \Sigma^{\mathcal{T}} \times \mathcal{Z}^{\mathcal{T}} \rightarrow \mathcal{S}$ is a sufficient statistic for $\mathbf{v}$ (resp. for $\mathbf{v}$ on $Z \subset \mathcal{Z}^{\mathcal{T}}$ ) if there is some $\tilde{\mathbf{v}}: \mathcal{Z}^{\mathcal{T}} \times \mathcal{S} \rightarrow$ $\mathbb{R}$ s.t. we have $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)=\tilde{\mathbf{v}}^{s\left(\sigma^{\mathcal{T}}, h\right)}(h)$ for each $h \in \mathcal{Z}^{\mathcal{T}}($ resp. $h \in Z)$ and $\sigma^{\mathcal{T}}$.

Proposition 4.27 (Sufficient statistics for optimal value functions). Let $\mathcal{T}$ be a trunk and $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$.
(1) The (joint) reach probabilities $\left(P^{\sigma^{\tau}}(h)\right)_{h \in \mathcal{Z} \mathcal{T}}$ provide a sufficient statistic for computing some reachably-optimal $\mathbf{v}$.
(2) The factored reach probabilities $\left(P_{1}^{\sigma^{\mathcal{T}}}(h), P_{2}^{\sigma^{\tau}}(h)\right)_{h \in \mathcal{Z} \mathcal{T}}$ provide a sufficient statistic for computing some counterfactually-optimal $\mathbf{v}$.
In particular, it suffices to keep $\left(P_{1}^{\sigma^{\mathcal{T}}}(I)\right)_{\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{1}}$ and $\left(P_{2}^{\sigma^{\mathcal{T}}}(I)\right)_{\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{2}}$.
Moreover, the decomposition of $\mathcal{Z}^{\mathcal{T}}$ into public states allows us to compute value functions in a more localized manner, by looking at a single subgame at a time. Formally, we have the following (well-known but previously unpublished) result:

Proposition 4.28 (Localization by public states). For any public state $S \subset \mathcal{Z}^{\mathcal{T}}$ :
(i) Both $\left(P^{\sigma^{\top}}(h)\right)_{h \in S}$ and $\left(P^{\sigma^{\top}}(h \mid S)\right)_{h \in S}$ provide a sufficient statistic for computing some reachably-optimal value function $\mathbf{v}^{\sigma^{\top}}(\cdot)$ on $S$.
(ii) $\left(P_{1}^{\sigma^{\mathcal{T}}}(I)\right)_{S \supset I \in \mathcal{I}_{1}}$ and $\left(P_{2}^{\sigma^{\mathcal{T}}}(I)\right)_{S \supset I \in \mathcal{I}_{2}}$ together provide a sufficient statistic for computing some counterfactually-optimal $\mathbf{v}^{\sigma^{\top}}(\cdot)$ on $S$.

We give an abstract proof of this result, showing that the definitions of appropriate optimal value functions can be rephrased to only depend on the statistics listed in the corresponding case of this proposition. This implies that it must be possible to compute the desired value function using those statistics only. Inspired by the poker literature, we refer to these statistics as range:

Definition 4.29 (Range). For a (possibly partial) strategy $\sigma$ and public state $S$, the range at $S$ (corresponding to $\sigma$ ) is defined as

$$
\begin{equation*}
\operatorname{rng}^{\sigma}(S):=\left(\left(P_{p}^{\sigma}(I)\right)_{S \supset I \in \mathcal{I}_{p}}\right)_{p=1,2} \tag{4.7}
\end{equation*}
$$

Similarly, we can talk about joint range $\left(P^{\sigma}(h)\right)_{h \in S}$ and normalized joint range $\left(P^{\sigma}(h \mid S)\right)_{h \in S}$. Moreover, [15] recently introduced the related notion of a public belief state, which is essentially the pair formed by a public state and a range at that state.

The following theorem shows that considering ranges at public states is not an ad-hoc choice, but in fact the canonical one. It states that changing the reach probability of a single history might change the value of any other history within the same public state. (Formally, the result uses common-knowledge public states, which [22] defines as the smallest subsets of $\mathcal{H}$ closed under infosets of both players.)
Theorem 5 (Public state minimality). Let $\mathcal{T}$ be a trunk, $\mathcal{Z}^{\mathcal{T}}$ its leaves, $S \subset$ $\mathcal{Z}^{\mathcal{T}}$ a common-knowledge public state, and $h_{0}, g \in S$. Suppose that trunk strategies $\sigma^{\mathcal{T}}$ and $\mu^{\mathcal{T}}$ render the same non-zero reach probabilities at $S$, except that $P^{\sigma^{\mathcal{T}}}(g) \neq P^{\mu^{\mathcal{T}}}(g)=0$.

Then there exists some game $G$, s.t. $\mathcal{T}$ is a trunk in $G$ and $\mathcal{Z}^{\mathcal{T}}$ the corresponding leaves, for which both $G\left(\sigma^{\mathcal{T}}\right)$ and $G\left(\mu^{\mathcal{T}}\right)$ each have a unique Nash equilibrium, $\sigma$ and $\mu$, and these satisfy $v_{1}^{\sigma}\left(h_{0}\right) \neq v_{1}^{\mu}\left(h_{0}\right)$.

While this result doesn't rule out more compact sufficient statistics for specific games or trunk strategies, it does say that any trunk is part of some game where anything less than data over whole public states will fail to be sufficient for many strategies. This yields the following corollary:

Corollary 4.30. Ranges over anything less than (common-knowledge) public states might fail to be a sufficient statistic for computing (reachably-optimal or better) value functions.

Output. With the results obtained so far, we see that many value functions can be evaluated independently in each public state and values of all histories within each $S$ can be computed all at once. In other words, we can view them as mappings of the type

$$
\mathbf{v}:(\text { public state } S, \text { range at } S) \mapsto \text { vector indexed by } h \in S
$$

However, public states typically have many more histories than information states (roughly quadratically), so the vector $\left(\mathbf{v}^{\text {range }}(h)\right)_{h \in S}$ will have a large
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-36.jpg?height=389&width=855&top_left_y=429&top_left_x=624)

Figure 2: An example of a game from Theorem 5
dimension. Since many algorithms run on a per-infoset basis anyway (e.g., CFR does), we can aggregate the individual history values into infoset values $\mathbf{V}_{p}^{(\cdot)}(I)$, and view value functions as mappings

$$
\mathbf{V}:(\text { public state } S, \text { range at } S) \mapsto \text { vector indexed by } I \in \mathcal{I}, I \subset S,
$$

where $\mathbf{V}^{(\cdot)}(I)=\sum_{h \in I} P^{(\cdot)}(h \mid I) \mathbf{v}_{p}^{(\cdot)}(h)$. Depending on the specific application, we might replace the output by the format that is the most appropriate (e.g., counterfactual values of infosets for DL-CFR). Since the resulting vectors will have a much lower dimension than the history-based representation, representing value functions in this manner makes them easier to approximate (e.g., by using a neural network).

Remark 4.31 (Recovering values of histories). One might argue that some algorithms are typically implemented over histories, so infoset-based value functions might fail to preserve enough information to run such algorithms. However, many such algorithms in fact only depend values of infosets, rather than histories. We can thus take the infoset values, $\mathbf{V}_{p}^{\sigma^{\tau}}(I)$, and decompose them into some "fake" history values $\left(\hat{\mathbf{v}}^{\sigma^{\tau}}(h)\right)_{h \in I}$. As long as we ensure that $\sum_{h \in I} P^{\sigma^{\tau}}(h \mid I) \hat{\mathbf{v}}^{(\cdot)}(h)=$ $\mathbf{V}_{p}^{\sigma^{\mathcal{T}}}(I)$ holds for all infosets $\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{p}$, these fake values will result in identical infoset values for $p$ in the whole trunk (Theorem1). This condition can be satisfied trivially by setting $\hat{\mathbf{v}}^{(\cdot)}\left(h_{0}\right):=P^{\sigma^{\mathcal{T}}}\left(h_{0} \mid I\right)^{-1} \mathbf{V}_{p}^{\sigma^{\top}}(I)$ for an arbitrary cf. reachable $h_{0} \in I$ and zeroing out the rest. This is a useful trick that allows us to combine infoset-based value functions with history-based implementations (in our case, depth-limited CFR).

Remark 4.32 (Values of public belief state). [15, Theorem 1] shows that infoset values coincide with supergradients of values of public belief states (which can be defined as $\left.v_{p}^{\sigma}(S, r)=\sum_{S \supset I \in \mathcal{I}_{p}} P^{\sigma}(I \mid S) \mathbf{v}_{p}^{\sigma}(I)\right)$. This means that instead of computing values for all infosets for a single range, we also have the option to only compute a single value for the whole public state but do so for various different ranges, such that we can approximate the partial derivatives with respect to reach probabilities of individual infosets (as these correspond to infoset values).

Taken to the extreme, this result can be extended to an equivalence between values of individual infosets and a single value for the whole bottom part of the game. (This is because the proof of Theorem 1 in [15] does not rely on the minimality of public states and, therefore, goes through even when the public state spans the whole width of the game tree.) In this paper, we chose to use the one-value-per-infoset representation because it seemed more suitable for discussing the differences between different types of optimality.

## 5. Value Functions in Partially-Observable Stochastic Games

In this section, we summarize the connection between the EFG model and the POSG model (traditionally used in multiagent RL). In particular, we observe that all of our results apply to POSGs as well, and that some of the existing POSG results are relevant (or even directly applicable) to value-functions in EFGs. We also list the most relevant results from the POSG literature.

Recall that the standard MARL model for zero-sum ${ }^{36}$ games is a partiallyobservable stochastic game (POSG), where players take actions which causes them to transition to a new state, obtain some reward, and receive some observation [6]. Another important difference between EFGs and POSGs is that while EFGs are tree-structured, POSGs can contain cycles.

Reasoning about general POSGs is difficult because we soon encounter nested beliefs - each player needing to maintain a belief over the current state of the game, a belief over the opponent's beliefs, a belief over the opponent's beliefs over the player's beliefs, and so on [24]. One way of countering these difficulties is to work with a restricted class of POSGs. This approach is certainly viable, and there are several works that analyze value functions in these settings [25, 24, 26, 27. However, nested beliefs are inherent to many domains of interest (including poker), so we also need some approach for tackling them. This is where the connection between POSGs and EFGs becomes relevant.

As shown in [23, 7, EFGs and POSGs should not be viewed as two unrelated models. Instead, EFGs are derived objects that can be obtained by "unrolling" some underlying POSG, such that every POSG $G$ has a tree-structured EFG representation $E$. While $E$ might be larger than the underlying POSG, it avoids nested beliefs by allowing players to use strategies defined on infosets (of which there are only finitely many) rather than on the difficult-to-handle nested beliefs. As a result, many concepts used for analyzing POSGs, such as policies, their expected rewards, and solution concepts, are often defined on the EFG representation ${ }^{37}$ of $G$ rather than on $G$ itself [28]. Most importantly for

[^19]our work, this is also true of value functions [5, 29, 30, which means that all of the present paper's results also apply to POSGs.

In the opposite direction, some of the results from the POSG literature are also highly relevant for value functions in EFGs [31, 32. A particularly relevant work is [5] which, translated into our terminology, studies the depth-limited expected-utility function $\left(\sigma_{1}^{\mathcal{T}}, \sigma_{2}^{\mathcal{T}}\right) \mapsto u_{p}^{\mathbf{v}}\left(\sigma_{1}^{\mathcal{T}}, \sigma_{2}^{\mathcal{T}}\right)$ (which the paper refers to as value function). It shows that for any reachably-optimal value function $\mathbf{v}$, $u_{p}^{\mathbf{v}}(\cdot, \cdot)$ is concave-convex and has a unique max-min value.

Finally, one downside of POSGs is that, by default, an EFG representation of a $P O S G$ does not define (non-trivial) public states, which are essential for localized computation of value functions (Proposition 4.28. Theorem 5). To reap the full benefit of the EFG representation, we can start with a minor extension of the POSG model, called factored-observation stochastic games (FOSGs) [7]. Indeed, each observation in a FOSG $G$ is explicitly decomposed into a private part and a public part, which naturally endows the corresponding EFG $E$ with a public partition - each public state $S$ in $E$ can be identified with a sequence $\vec{o}_{\text {pub }}$ of public observations in $G$. Similarly, each information set can be identified with a pair $\left(\vec{o}_{\text {pub }}, \vec{o}_{\text {priv }(p)}\right)$, where $\vec{o}_{\operatorname{priv}(p)}$ is a sequence of private observations (and actions) of player $p$.

An important benefit of the FOSG-EFG correspondence is that identifying public states and information sets with sequences of observations gives them additional structure. To the extent that similar game states (in terms of what the players observe) have similar values, this structure might allow for better generalization when approximating value functions by neural networks. For this reason, all games in our experiments (Section 6) are initially specified using the FOSG model.

## 6. Experimental Evaluation

The goal of this section is to complement the theoretical contribution (Sections 4 and 5 by an analysis of the practical side of value functions - approximating value functions by value networks and using them for depth-limited solving. In particular, we train a value network that would, in the limit of infinite resources, act as a counterfactually optimal value function (Definition 4.15) and use it in conjunction with a depth-limited version of CFR (Example 4.12). We call the resulting algorithm $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$, and describe it in detail in Section 6.1. From results such as [1, 9, 33, 3, we already know that approaches based on CFR and value functions can achieve impressive results. Instead of showing that the algorithm can solve large domains, our experiments thus focus on medium-sized domains where it is possible to thoroughly investigate the algorithm's behaviour and compute exact exploitabilities. We consider three domains with qualitatively different properties: Leduc hold'em poker (LH), a standard benchmark for imperfect information games, and (similarly-sized variants of) imperfect information goofspiel (GS) and imperfect information oshi-zumo (OZ).

We start by describing the experimental setup and a universal encoding that avoids the need for hand-crafting domain features (Section 6.1). Afterwards, we
investigate several questions: First of all, since most past results only evaluated CFR and its variants on poker and poker-like domains (which have quite specific properties) [34, we should verify whether the depth-limited CFR approach empirically works beyond poker. This turns out to be the case - as a byproduct of our other experiments, we see that it is possible to train value networks that achieve low validation loss in all three domains and that coupling these value networks with $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$produces strategies with low exploitability (Section 6.2). Second, we study the impact of the value network's quality on the strength of the resulting depth-limited solver, concluding that $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$with a near-optimal value network successfully finds an $\epsilon$-Nash equilibrium strategy (Section 6.2). In particular, validation Huber loss of 0.001 was sufficient to achieve exploitability below 0.01 . Third, we look at how well the network generalizes to previously unseen situations (Section 6.3). We observe that $\mathrm{DL}_{\mathrm{L}} \mathrm{CFR}_{\mathrm{NN}}^{+}$ achieves a low exploitability despite encountering inputs significantly different from those present in the training data. Moreover, we show that the value network generalizes to unseen public states when there is enough public states to generalize from, suggesting that the FOSG encoding is suitable for value networks. Finally, since $\mathrm{DL}^{-\mathrm{CFR}} \mathrm{NN}^{+}$with a perfect value network would behave identically to CFR-D [17], we qualitatively compare the two algorithms (Section 6.4). We see that the value network achieves low loss on the inputs encountered in CFR-D and that both algorithms produce similar strategies on the tested domains.

In summary, our experiments demonstrate that the value function approach to depth-limited solving is viable in domains other than poker and that value functions can be approximated using a universal, domain-independent encoding and architecture.

### 6.1. Experimental Setup

We now describe the algorithm used for depth-limited solving, its value network, the experimental domains, the process used for generating the corresponding training data, and the encoding of this data.

### 6.1.1. Depth-limited Solver

We use DL-CFR ${ }_{\mathrm{NN}}^{+}$- a depth-limited version of CFR + (Example 4.12, 35]) coupled with a value network. Each iteration, the algorithm computes the new strategy in the trunk and collects the ranges for all public states at the depth limit. Each public state $S$ and range is encoded (in a way described below) and used as a query to the value network, which outputs a vector of values, one for each infoset in $S$. Since our implementation of CFR+ runs on histories rather than on infosets, we then use the simple trick from Remark 4.31 to convert the infoset values to history values. These values are then back-propagated by CFR+ as usual. When using DL-CFR ${ }_{\mathrm{NN}}^{+}$, we always run it for 1000 iterations. The implementation of $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$is available online at 36.

### 6.1.2. Training Data

The training data is produced by the following process. First, we generate a random trunk strategy at each infoset in the trunk by choosing either a fully
mixed strategy drawn from a uniform probability distribution $90 \%$ of the time or a pure strategy for a randomly chosen action $10 \%$ of the time. We then fix this trunk strategy and solve the bottom of the game using 1000 iterations of CFR+. In Section 4.2.3, we referred to this technique as value solving and showed that it produces counterfactually (near-)optimal values, so $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$should with a perfect network - converge to a Nash equilibrium (by Proposition 4.20). Finally, we create a single training data point out of each public state $S$ at the depth-limit - the input is formed by $S$ and the corresponding range at $S$ (Definition 4.29), the output is formed by the vector of counterfactual infosetvalues (Section 3.3). (We explain this encoding in more detail, but only after describing the experimental domains.) As a result, each trunk strategy yields as many training data points as there are infosets at the depth limit (see Table 3). We use $90 \%$ of the data for training and $10 \%$ for validation.

Since obtaining training data can be costly, we investigated what is the minimum amount of training data needed for satisfactory performance on unseen inputs. We used training sets of increasing size from 1000 up 50000 (step size $1000)$ depending on the domain, and computed the validation loss on a fixed data set. In all three domains, between below 10000 samples were sufficient for achieving a validation Huber loss of 0.001 . We however opted to use a larger amount of data to run our main experiments on to analyze DL-CFR ${ }_{\mathrm{NN}}^{+}$using a well-trained value network (exact numbers in Table 22). Note that the necessary amount of data might grow significantly with increasing domain size (for example, DeepStack [1 used 10 million random trunk strategies). In such larger domains, it might, therefore, be advantageous to use a combination of (a) a more involved training process (such as the one used in [15]) that puts more weight on inputs that are more likely to be relevant for the algorithm and (b) an approach like continual resolving [1, 22], which splits the game at multiple depths.

### 6.1.3. Value Network

All of the networks use a standard fully-connected feed-forward network, rectified linear units for the hidden layers and linear activation for the output layer, 3-6 layers, and layer-width 5-10x the input size (see Table 24). They were trained for 1000 epochs by using the Adam optimizer [37] while minimizing Huber loss. The reasoning behind the choice of loss function is further explored in Section 6.2, Before settling on these choices, we experimented with various other options and hyperparameter values; for more details, see Appendix D and E. 1

### 6.1.4. Experimental Domains

We now give a high-level description of the domains used in the experiments. For a more detailed description, see Appendix C. The first domain we use is Leduc hold'em (LH) 38 - a standard small variant of poker with a deck of six cards, two rounds of betting, and fixed bet sizes. Since poker is an imperfect information game with rather specific properties [34, we also use two other domains which are qualitatively different from it (and from each other). One of these games is imperfect information oshi-zumo (OZ) with board size 3 and a

| Domain | Layers | Neurons | Training Strategies | Tr. Datapoints |
| :---: | :---: | :---: | :---: | :---: |
| GS | 5 | 500 | 2000 | 18000 |
| OZ | 4 | 400 | 2000 | 34000 |
| LH | 6 | 200 | 812 | 70644 |

Table 2: The final value network architectures and the amounts of training data. Recall that each trunk strategy yields as many training samples as there are public states at the depth limit.

| Game Property | Goofspiel | Oshi-zumo | Leduc hold'em |
| :---: | :---: | :---: | :---: |
| private actions | Yes | Yes | No |
| early terminal nodes | No | Yes | Yes |
| constant-size infosets | No | No | Yes |
| rounds | 5 | 8 | 4 |
| total infosets | 5000 | 1600 | 4000 |
| total histories | 56000 | 20000 | 61000 |
| depth-limit (rounds) | 3 | 4 | 2 |
| publ. states at depth-limit | 9 | 17 | 87 |

Table 3: Properties of the domains used for empirical evaluation.
budget size 8. In this game, each player controls a wrestler and spends some amount of energy (technically called "coins") each round in an attempt to pushing the opponent towards the edge of the board. In the imperfect information variant of the game, the players learn which of them spent more each round, but not what the exact amount was. The last domain we use is a 5 -card variant of imperfect information goofspiel (GS). In this game, a deck of cards $\{1, \ldots, 5\}$ is auctioned off card by card, with each player trying to maximize the sum of card-values they win. However, instead of money, each player has their own deck $\{1, \ldots, 5\}$ that they for betting. The different properties of all three domains are summarized in Table 3. For the purpose of depth-limited solving, we split each of the domains in a trunk and bottom after half of the rounds has passed. When calculating exploitability, we normalize all utilities to make the results comparable.

### 6.1.5. Domain-independent Encoding

On the mathematical level, each input to the network is a pair, a public state $S$ and a range $r$ at $S$ (a vector of reach probabilities, one for each infoset at $S$; Definition 4.29 and the corresponding output is a vector of counterfactual values (with the same shape as $r$ ). We now explain how $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$encodes these objects when communicating with its value network. The FOSG representation identifies each public state $S$ with a sequence of public observations $\vec{o}_{\text {pub }}=$ $\left(o_{1}, \ldots, o_{n}\right)$, where each $o_{i}$ belongs to some set of possible observations. Since this set is discrete and rather small in all of our domains (and in many others), we can enumerate it, use a one-hot vector to encode each $o_{i}$, and concatenate
these vectors to obtain an encoding of the full sequence $\left(o_{1}, \ldots, o_{n}\right)$. For example, there are three possible public outcomes of each round of goofspiel, depending on whether the currently-auctioned card is given either to player one, player two, or discarded (when the players make identical bids), so the corresponding observations would be encoded as $[1,0,0],[0,1,0]$, resp. $[0,0,1]$. If player one wins the first round, loses the second round, and the players are currently bidding for the third card, the corresponding public state would be encoded as $[1,0,0,0,0,1]$.

The encoding of ranges also relies on the FOSG representation. In a FOSG, each infoset $I \in \mathcal{I}_{p}$ is identified with a pair $\left(\vec{o}_{\text {pub }}, \vec{o}_{\text {priv }(p)}\right)$, where $\vec{o}_{\text {pub }}$ corresponds to the public set $S$ that contains $I$ and $\vec{o}_{\operatorname{priv}(p)}$ is some sequence of $p$ 's private observations and actions. Since the publicly observable actions (e.g., bets in Leduc hold'em) are already a part of $\vec{o}_{\text {pub }}$, we can assume that $\vec{o}_{\text {priv }(p)}$ only contains $p$ 's private observations and their non-public actions (e.g., bid sizes in imperfect-information goofspiel and oshi-zumo). The key insight is that in the context of a specific public state, we only need the private sequence $\vec{o}_{\operatorname{priv}(p)}$ to specify an infoset. Consequently, we construct a universal encoding that works all public states, by fixing some enumeration of the set of all sequences $\vec{o}_{\mathrm{priv}(\mathrm{p})}$ that can appear in some public state at the depth limit. (For example, in our variant of goofspiel, there are no private observations but all actions - i.e., using up one of the cards $\{1, \ldots, 5\}$ - are private. The set of all possible private sequences $\vec{o}_{\operatorname{priv}(p)}$ at the start of the third round can, therefore, be identified with $\left\{\left(c_{1}, c_{2}\right) \mid c_{i}=1, \ldots, 5, c_{1} \neq c_{2}\right\}$.) With this enumeration, the $p$ 's part of a range $r$ at a public state $S$ can be encoded as a vector $\left(r_{1}, \ldots r_{n}\right)$, where $r_{i}$ is the reach probability of the infoset that corresponds to the $i$-th private sequence $\vec{o}_{\text {priv }(p)}$. If this $\vec{o}_{\operatorname{priv}(p)}$ does not correspond to any infoset at the given public state, we set $r_{i}$ to zero. (For example, our goofspiel scenario $S$ where player one wins the first round and looses the second is incompatible with any sequence $\vec{o}_{\text {priv }(p)}$ where they bid $1-$ smallest possible amount - in the first round or 5 largest possible amount - in the second round.) To encode the full range, we concatenate the part of player 1 with the part of player 2 . Since counterfactual values have the same structure as ranges, we use the same method for their encoding (or rather, decoding). Overall, a major contribution of this encoding is that it works for any domain described as a FOSG. Its main limitation (which we leave for future work) is that the width of the network's input layer scales with the number of private observation sequences. Fortunately, this number is manageable in the majority of games studied in the CFR literature so far - e.g., in all variants of poker including no-limit Texas hold'em [1], in liar's dice 39], and in all domains considered in this paper.

### 6.2. Value Network's Impact on Depth-limited Solving

While we ultimately care about the exploitability of $\mathrm{DL}_{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$'s strategy, computing it is costly (Proposition 4.5), so we would like to know whether the value network's quality is a reliable proxy for exploitability and if so, how low validation loss should we aim for.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-43.jpg?height=394&width=1213&top_left_y=426&top_left_x=448)

Figure 3: The relationships between the value network's Huber validation error ( x -axis) and on the exploitability of $\mathrm{DL}-\mathrm{CFR}_{\mathrm{NN}}^{+}$. The figures display aggregate data from training runs with 100 different initializations. We grouped the data into eight buckets based on the validation loss. For each bucket, we display the interval between minimum and maximum exploitability (the red line) and the average exploitability (black dot) achieved with the corresponding value networks. The green line indicates the exploitability of CFR $+(1000$ iterations, same as $\mathrm{DL}-\mathrm{CFR}_{\mathrm{NN}}^{+}$). The exploitabilities of $\mathrm{DL}^{+} \mathrm{CFR}_{\mathrm{NN}}^{+}$with a constant value function ( 0.5 in GS, 0.7 in $\mathrm{OZ}, 0.2$ in LH ) are not shown as they would not fit into the graph.

Before analyzing the main question, we first look at a one aspect value network training that hasn't been addressed by previous work: which loss function should we use for training and validation (the main candidates being Huber, $l_{1}$, and $l_{\infty}$ [1]). The details are presented in Appendix D. In summary, for each domain and loss function combination, we trained value networks to minimize the given loss while calculating the validation loss for all three loss functions. We observed that the distinction between the losses is not significant. We therefore choose Huber as it has been used in prior work [1]. We also performed an in-depth analysis of custom loss functions (Appendix D) before concluding that standard losses are sufficient. As a result, we train all value networks by minimizing the Huber loss while measuring all three losses (Huber, $l_{1}, l_{\infty}$ ) on validation data and reporting the one that seems the most intuitive for the given experiment.

To investigate the relationship between network loss and exploitability, we ran 100 training runs of the value network (with different initializations) and continually saved checkpoints of the current weights and the corresponding validation losses. We used each checkpoint for $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$(1000 iterations), measured the exploitability of the resulting strategy, and plotted the aggregate data in Figure 3. As a baseline, we computed the exploitability of the standard CFR + after 1000 iterations. As a sanity check, we measured the exploitability of $\mathrm{DL}-\mathrm{CFR}_{\mathrm{NN}}^{+}$with a constant value function (that ignores the input and predicts value 0 for every infoset) - this always resulted in high exploitability ( 0.5 in GS, 0.7 in OZ, 0.2 in LH). In all cases, we see a strong connection between validation loss and the resulting exploitability. The networks are able to reach training Huber loss below 0.001 and exploitability below 0.01 - in other words, $\mathrm{DL}-\mathrm{CFR}_{\mathrm{NN}}^{+}$performs on par with CFR+ that runs on the full tree.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-44.jpg?height=389&width=1194&top_left_y=429&top_left_x=460)

Figure 4: Qualitative comparison of CFR-D, DL-CFR ${ }_{\mathrm{NN}}^{+}$, and training data in a specific public state in Leduc hold'em. Each cell corresponds to a reach probability (resp. counterfactual value) of one infoset in the given iteration. The values in the left column depict the values computed by CFR-D's resolving, $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$value network's predictions are in the middle, and the training data points closest to the $\mathrm{DL}_{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$data are on the right.

### 6.3. Generalizing to Unseen Inputs

To ensure that the value network isn't merely memorizing all data, we investigate its ability to generalize to unseen situations. Recall that since the validation losses are low, the network must be able to do some amount of generalization. Moreover, as shown in Appendix E.2 the validation error remains low (in the sense of corresponding to an acceptable exploitability according to Figure 3) even when we validate on inputs that would appear in CFR-D (i.e., on inputs $\mathrm{DL}_{\mathrm{L}}-\mathrm{CFR}_{\mathrm{NN}}^{+}$would ask for if its value network was perfect). However, we can also look at the network's generalization ability in more detail.

First, we investigate how different the training ranges (Definition 4.29) are from those used by $\mathrm{DL}^{-\mathrm{CFR}_{\mathrm{NN}}^{+}}$. To do so, we look at a single run of the algorithm in each domain. Focusing on a single public state, we look at each input requested by $\mathrm{DL}_{-\mathrm{CFR}_{\mathrm{NN}}}^{+}$and compare it with the training data point that is closest to it in terms of Euclidean distance (see the middle and right columns of Figure 4, resp. Appendix E.3, for a visualization of the data in LH, resp. GS and OZ). In all cases, we see that even the closest training data points are significantly different from those requested by $\mathrm{DL}_{\mathrm{L}} \mathrm{CFR}_{\mathrm{NN}}^{+}$.

Second, we evaluate the value network's ability to generalize to parts of the game not encountered in the training. We do this by withholding all data about a single public state from the training data, training the network on the remaining data, and computing validation loss on the withheld data. We repeat this calculation for all public states at the depth limit, and report the results for Leduc hold'em in Figure 5 and for the other domains in Appendix E.4. The generalization is very poor in goofspiel, reasonable in oshi-zumo, and very strong in Leduc hold'em. Since these domains have 9, 17, resp. 87 public states at the depth limit, this is suggests (though isn't a hard proof) that the encoding is suitable for value networks - at least in our domains, it allows for generalization to the extent that there is enough data to generalize from.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-45.jpg?height=494&width=952&top_left_y=433&top_left_x=565)

Figure 5: We show $l_{\infty}$ validation errors over withheld public states in Leduc hold'em. Each orange bar corresponds to the $l_{\infty}$ validation loss of a value network trained on all data except for the public state $S$ with given index. For comparison, we show the loss (at $S$ ) of the value network trained on all data (green) and the average loss (at $S$ ) of 10 randomly initialized value networks.

### 6.4. Qualitative Comparison to $C F R-D$

With a perfect value network, $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$would behave identically to CFRD. As our last experiment, we investigate whether this similarity also appears in practice. We run CFR-D on all three domains, record the corresponding ranges of public states at the bottom of the trunk, and compare them to ranges that appear in $\mathrm{DL}^{-\mathrm{CFR}_{\mathrm{NN}}}$. Additionally, we show the corresponding values (subgame solutions) which have been computed by both algorithms in the bottom row. The results for Leduc hold'em poker are depicted in the left and middle columns of Figure 4 , the data for the other two domains can be found in Appendix E. 3 . We see that - at least in all of our domains - the two algorithms encounter very similar ranges and converge to similar strategies.

## 7. Discussion

In this section, we summarize the most-relevant conceptual and theoretical results from related work and present our conclusions.

### 7.1. Related Work

The two results that inspired the paper the most are CFR-D [17] and DeepStack [1]. Roughly speaking, a portion of this paper might be viewed as a domain- and algorithm- independent extension of CFR-D and DeepStack (minus the parts about continual resolving and evaluating on a large domain):

In CFR-D, the letter D stands for "decomposition". However, the algorithm can be understood as a form of depth-limited solving, first that we are aware of in imperfect-information games. The key differences between CFR-D and the present paper are: (1) CFR-D always solves the trunk using CFR. (2) Upon reaching the bottom of the trunk, CFR-D solves the subgames by CFR or some
other game-solving algorithm (instead of evaluating them using a value function). (3) CFR-D is only evaluated on Leduc hold'em. The paper uses the notion of trunk and subgames but does not define them formally. It also points out that for CFR-D to work, the subgame-strategies need to be mutual counterfactual best responses, not merely mutual best responses (i.e., Nash equilibria). This observation inspired us to differentiate between reachable and counterfactual optimality and prompted our search for (counter-) Example 4.18 .

DeepStack has three relevant contributions: (1) it builds on top of CFR-D, but rather than solving subgames by CFR, it trains a value network to imitate CFR's behaviour on subgames. (2) It introduces continual resolving - iteratively applying depth-limited solving in online play, and thus splitting the game on multiple depth levels. (3) It demonstrates that with (1), (2), and certain pokerspecific tricks, it is possible to achieve superhuman performance in full-sized poker (two-player no-limit Texas hold 'em). The key limitation of DeepStack is that both the value-network architecture and the overall algorithm is specific to poker and CFR. Additionally, due to imitating CFR-D, the value-network's training target is only defined in terms of CFR.

A second related line of work, which seemingly does not have anything to do with value functions, revolves around multivalued-states [4]: Where [1] replaces each trunk-leaf history by a call to a value function, 4] replaces each trunk-leaf by a single decision of which strategy to use for the remainder of the game. In Section 4.3.1, we show that multivalued states can be understood as a computationally-convenient implementation of a particular type of an optimal value function (Theorem 4).

A third line of closely related work is the study of value functions in partiallyobservable stochastic games, discussed further in Section 5. Two results worth highlighting here are: (1) By [7, 23], the value functions studied in the present paper also apply to POSGs. (2) The paper [5] proves that (what we would call) reachably-optimal values aggregated over the whole bottom of the trunk are concave-convex and have a unique max-min value.

In summary, the value-function approach [1] and the multivalued-states approach [4] are two ways of obtaining a similar result and can likely be combined. Both of these can be applied to POSGs (and FOSGs [7]), and some of the POSG results have implications for both approaches.

There are also several recent relevant results (which only came out after this paper has been written, or around the same time). The papers [40, 41] extend alpha-beta pruning to imperfect-information games, and some of its definitions are relevant to our Section 4.1. The paper introducing ReBeL [15] shows that infoset-values can be viewed as derivatives of values aggregated over public states and describes an online value-network training procedure similar to AlphaZero [42]. Similarly, Player of Games [9] is an extension of AlphaZero and DeepStack [1] which performs competitively in both perfect- and imperfect-information games (chess, Go, poker, and stratego).

Finally, as the terminology suggests, the values considered in this paper are connected to the $V$ - and $Q$-values used in (multiagent) (reinforcement) learning (MARL). However, since both our setting and typical MARL settings need to
be significantly simplified to overlap, the results in the two lines of work do not inform each other as much. Indeed, the frequent differences are: (a) MARL often assumes general-sum rewards or $N \geq 3$ players, (b) we assume a finite-horizon, no discount factor, and rewards only in terminal states, (c) MARL often assumes simultaneous moves but otherwise perfect information (i.e., Markov games 43]), and hence uses stationary strategies. However, where the two settings intersect, the defined value functions coincide. For example, in MDPs (with a finite horizon, no discount factor, rewards in terminal states only, and a tree-structured statespace), universally-optimal value functions coincide with $v^{*}(s)$ [44. Similarly, universally-optimal value functions in two-player zero-sum Markov games (with the same qualifications) coincide with the minimax state values 43 and can be used to define the $Q$-functions used in [45, 46].

### 7.2. Conclusion

In this paper, we gave an accessible description of basic notions used in the CFR literature (Section 3) and introduced a number of concepts that enable reasoning about depth-limited games and value functions (Section 4.1). We proved that different degree of value-function optimality is required for different calculations and explained how to obtain the key types of value functions (Sections 4.2.2 4.2.4). Additionally, we proved that public belief states provide the necessary and sufficient context for computing value functions (Theorem 5). Our description allows viewing Deepstack's value functions 1 and Brown et al.'s multivalued states [4] as two instances of a single unifying framework (Theorem 4). Moreover, the results also apply to partially-observable stochastic games and their recent extension, factored-observation stochastic games [7] (Section 5). The theory shows that depth-limited solving is applicable to arbitrary domains and various algorithms. However, due to its success in recent years, our experimental evaluation focused on CFR. We showed that adopting the FOSG formalism allows for a simple domain-independent encoding which can be used for input and output of a value function (Section 6.1.5). We demonstrated the suitability of this encoding by showing that the resulting value network can generalize to unseen public states (Section 6.3). In three distinct domains, we used this encoding to train a simple feed-forward neural network that approximates an optimal value function. We then implemented a depth-limited version of CFR that utilizes this network. We performed an extensive experimental evaluation of this setup. Most importantly, we confirmed that the value network's error is strongly correlated with the exploitability of the strategy found by the corresponding DL-CFR ${ }_{\mathrm{NN}}^{+}$ (Figure 3), achieving performance that is as good as that of CFR with access to the full game.

Overall, we have shown that depth-limited solving is a viable and robust option for a range of imperfect-information games beyond poker.

## Acknowledgements

This work was both supported by the Czech science foundation grant no. 18-27483Y and RCI grant CZ.02.1.01/0.0/0.0/16 019/0000765.

## References

[1] M. Moravčík, M. Schmid, N. Burch, V. Lisý, D. Morrill, N. Bard, T. Davis, K. Waugh, M. Johanson, M. Bowling, Deepstack: Expert-level artificial intelligence in heads-up no-limit poker, Science 356 (6337) (2017) 508-513.
[2] N. Brown, T. Sandholm, Superhuman AI for heads-up no-limit poker: Libratus beats top professionals, Science (2017) eaao1733.
[3] J. Serrino, M. Kleiman-Weiner, D. C. Parkes, J. B. Tenenbaum, Finding friend and foe in multi-agent games, arXiv preprint arXiv:1906.02330 (2019).
[4] N. Brown, T. Sandholm, B. Amos, Depth-limited solving for imperfect-information games, arXiv preprint arXiv:1805.08195 (2018).
[5] A. J. Wiggers, F. A. Oliehoek, D. M. Roijers, Structure in the value function of two-player zero-sum games of incomplete information, in: Proceedings of the Twenty-second European Conference on Artificial Intelligence, IOS Press, 2016, pp. 1628-1629.
[6] E. A. Hansen, D. S. Bernstein, S. Zilberstein, Dynamic programming for partially observable stochastic games, in: AAAI, Vol. 4, 2004, pp. 709-715.
[7] V. Kovařík, M. Schmid, N. Burch, M. Bowling, V. Lisý, Rethinking formal models of partially observable multiagent decision making, Artificial Intelligence (2021) 103645.
[8] N. Brown, T. Sandholm, Superhuman AI for multiplayer poker, Science 365 (6456) (2019) 885-890.
[9] M. Schmid, M. Moravčík, N. Burch, R. Kadlec, J. Davidson, K. Waugh, N. Bard, F. Timbers, M. Lanctot, Z. Holland, et al., Player of games, arXiv preprint arXiv:2112.03178 (2021).
[10] A. Celli, A. Marchesi, T. Bianchi, N. Gatti, Learning to correlate in multi-player general-sum sequential games, Advances in Neural Information Processing Systems 32 (2019) 13076-13086.
[11] S. K. Jakobsen, T. B. Sørensen, V. Conitzer, Timeability of extensive-form games, in: Proceedings of the 2016 ACM Conference on Innovations in Theoretical Computer Science, 2016, pp. 191-199.
[12] J. Y. Halpern, R. Pass, Sequential equilibrium in games of imperfect recall., in: KR, 2016, pp. 278-287.
[13] R. Selten, Reexamination of the perfectness concept for equilibrium points in extensive games, Economics (1974).
[14] M. Zinkevich, M. Johanson, M. Bowling, C. Piccione, Regret minimization in games with incomplete information, in: Advances in neural information processing systems, 2008, pp. 1729-1736.
[15] N. Brown, A. Bakhtin, A. Lerer, Q. Gong, Combining deep reinforcement learning and search for imperfect-information games, arXiv preprint arXiv:2007.13544 (2020).
[16] P. Auer, N. Cesa-Bianchi, Y. Freund, R. E. Schapire, The nonstochastic multiarmed bandit problem, SIAM journal on computing 32 (1) (2002) 48-77.
[17] N. Burch, M. Johanson, M. Bowling, Solving imperfect information games using decomposition., in: AAAI, 2014, pp. 602-608.
[18] Y. Shoham, K. Leyton-Brown, Multiagent systems: Algorithmic, game-theoretic, and logical foundations, Cambridge University Press, 2008.
[19] V. Lisý, Alternative selection functions for information set Monte Carlo tree search, Acta Polytechnica 54 (5) (2014) 333-340.
[20] N. Burch, Time and space: Why imperfect information games are hard, Ph.D. thesis, University of Alberta (2017).
[21] J. Heinrich, M. Lanctot, D. Silver, Fictitious self-play in extensive-form games., in: ICML, 2015, pp. 805-813.
[22] M. Šustr, V. Kovařík, V. Lisý, Monte Carlo continual resolving for online strategy computation in imperfect information games, in: Proceedings of the 18th International Conference on Autonomous Agents and Multiagent Systems (AAMAS), 2019, pp. 224-232.
[23] F. Oliehoek, N. Vlassis, et al., Dec-POMDPs and extensive form games: equivalence of models and algorithms, Ias technical report IAS-UVA-06-02, University of Amsterdam, Intelligent Systems Lab, Amsterdam, The Netherlands (2006).
[24] M. Dermed, L. Charles, Value methods for efficiently solving stochastic games of complete and incomplete information, Ph.D. thesis, Georgia Institute of Technology (2013).
[25] H. L. Cole, N. Kocherlakota, Dynamic games with hidden actions and hidden states, Journal of Economic Theory 98 (1) (2001) 114-126.
[26] K. Horák, B. Bošanský, Solving partially observable stochastic games with public observations, in: Proceedings of the AAAI Conference on Artificial Intelligence, Vol. 33, 2019, pp. 2029-2036.
[27] K. Horák, B. Bošanský, M. Pěchouček, Heuristic search value iteration for one-sided partially observable stochastic games, in: Thirty-First AAAI Conference on Artificial Intelligence, 2017, pp. 558-564.
[28] F. A. Oliehoek, Decentralized POMDPs, in: Reinforcement Learning, Springer, 2012, pp. 471-503.
[29] O. Buffet, J. Dibangoye, A. Delage, A. Saffidine, V. Thomas, On Bellman's optimality principle for zs-POSGs, arXiv preprint arXiv:2006.16395 (2020).
[30] A. Delage, O. Buffet, J. Dibangoye, HSVI fo zs-POSGs using concavity, convexity and Lipschitz properties, arXiv preprint arXiv:2110.14529 (2021).
[31] F. A. Oliehoek, Sufficient plan-time statistics for decentralized POMDPs, in: Twenty-Third International Joint Conference on Artificial Intelligence, 2013, pp. 302-308.
[32] J. S. Dibangoye, C. Amato, O. Buffet, F. Charpillet, Optimally solving Dec-POMDPs as continuous-state MDPs, Journal of Artificial Intelligence Research 55 (2016) 443-497.
[33] R. Zarick, B. Pellegrino, N. Brown, C. Banister, Unlocking the potential of deep counterfactual value networks (2020). arXiv:2007.10442
[34] V. Kovařík, D. Milec, M. Šustr, D. Seitz, V. Lisý, Fast algorithms for poker require modelling it as a sequential bayesian game, arXiv preprint arXiv:2112.10890 (2021).
[35] O. Tammelin, CFR+, CoRR, abs/1407.5042 (2014).
[36] Gtlib2 implementation of DL-CFR-NN, https://gitlab.fel.cvut.cz/ game-theory-aic/GTLib2/-/blob/master/algorithms/cfr_dl.cpp, accessed: 2020-09-17 (2020).
[37] D. P. Kingma, J. Ba, Adam: A method for stochastic optimization, arXiv preprint arXiv:1412.6980 (2014).
[38] F. Southey, M. P. Bowling, B. Larson, C. Piccione, N. Burch, D. Billings, C. Rayner, Bayes' bluff: Opponent modelling in poker, arXiv preprint arXiv:1207.1411 (2012).
[39] N. Burch, M. Lanctot, D. Szafron, R. Gibson, Efficient Monte Carlo counterfactual regret minimization in games with many player actions, Advances in neural information processing systems 25 (2012).
[40] B. Zhang, T. Sandholm, Small nash equilibrium certificates in very large games, Advances in Neural Information Processing Systems 33 (2020) 7161-7172.
[41] B. Zhang, T. Sandholm, Finding and certifying (near-) optimal strategies in black-box extensive-form games, in: AAAI Workshop on Reinforcement Learning in Games, 2021, pp. 1-10.
[42] D. Silver, J. Schrittwieser, K. Simonyan, I. Antonoglou, A. Huang, A. Guez, T. Hubert, L. Baker, M. Lai, A. Bolton, et al., Mastering the game of Go without human knowledge, Nature 550 (7676) (2017) 354.
[43] M. L. Littman, Markov games as a framework for multi-agent reinforcement learning, in: Machine learning proceedings 1994, Elsevier, 1994, pp. 157-163.
[44] R. S. Sutton, A. G. Barto, Reinforcement learning: An introduction, MIT press, 2018.
[45] J. Hu, M. P. Wellman, Nash $q$-learning for general-sum stochastic games, Journal of machine learning research 4 (Nov) (2003) 1039-1069.
[46] A. Greenwald, K. Hall, R. Serrano, et al., Correlated Q-learning, in: ICML, Vol. 3, 2003, pp. 242-249.
[47] M. Šustr, M. Schmid, M. Moravčík, N. Burch, M. Lanctot, M. Bowling, Sound search in imperfect information games, arXiv preprint arXiv:2006.08740 (2020).
[48] M. Buro, Solving the oshi-zumo game, in: Advances in Computer Games, Springer, 2004, pp. 361-366.
[49] Loss functions for DL-CFR-NN, https://gitlab.fel.cvut.cz/seitzdom/value_func_repo_python/-/ blob/master/nn/loss_functions.py, accessed: 2020-10-20 (2020).

## A. Non-Uniqueness of Value Functions

Having explained how to compute and apply optimal value functions, it remains to ask about their uniqueness properties. Unfortunately, the answers we offer in this section are often negative: there might be value functions that are all optimal in the same sense, yet prescribe different values ${ }^{38}$ We will see that if one aims to use value functions in conjunction with function-approximation techniques, the natural follow-up question is whether the set of optimal value functions is at least convex. We were unable to establish a conclusive result in either direction, but we have found some evidence that the answer might be positive. The remainder of this section is devoted to giving further details and counterexamples on this topic.

Proposition A. 1 (Optimal values aren't unique). There exists $G$ and $\mathcal{T}$ for which there are multiple distinct universally-optimal value functions. This holds even after infoset-aggregation $\mathbf{V}_{p}^{\sigma^{\mathcal{T}}}(I):=\sum_{h \in I} P^{\sigma^{\mathcal{T}}}(h \mid I) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h)$.

Proof. The sequential representation of matching pennies demonstrates this well. Suppose the trunk $\mathcal{T}$ only consists of the root node where player 1 makes their decision and let $\sigma^{\mathcal{T}}$ be the uniform strategy. Given the uniform trunk strategy, any strategy that player 2 selects in the bottom of the game is going to extend $\sigma^{\mathcal{T}}$ optimally (even universallyoptimally). However, different player 2 strategies will result in different tradeoffs between $v_{p}^{(\cdot)}$ (heads) and
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-52.jpg?height=250&width=245&top_left_y=1157&top_left_x=1306) $v_{p}^{(\cdot)}$ (tails), thus proving the first part of the proposition. Since the singletons \{heads\} and \{tails\} each constitute a single infoset of player one, the infoset values are not unique either ${ }^{39}$

In light of Remark 4.32, this non-uniqueness should come as no surprise: 15, Theorem 1] establishes that infoset values correspond to supergradients of public belief state (PBS) values. But since the PBS value function is concave but not necessarily differentiable, the supergradients might not be unique. As a result, infoset values do not have to be unique either.

Some hope for infoset-value uniqueness stems from the fact that for a given trunk strategy $\sigma^{\mathcal{T}}$ and opponent strategy $\sigma_{-p} \supset \sigma_{-p}^{\mathcal{T}}$, all optimal responses $\sigma_{p} \supset \sigma_{p}^{\mathcal{T}}$ must lead to same values infoset values of $\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{p}$ (trivially, since if $V_{p}^{\sigma}(I)$ is equal to $\max _{\tilde{\sigma}_{p} \in \Sigma_{p}} V_{p}^{\tilde{\sigma}_{p}, \sigma_{-p}}(I)$, it can only depend on $\left.\sigma_{-p}\right)$. If one player has only a single optimal extension of the trunk strategy, the other player will thus have unique optimal infoset values.

[^20]Given that optimal values aren't unique, we might instead ask whether the corresponding set of functions is convex. To see why this would be relevant, suppose that one half of a training set corresponds to an optimal value function $\mathbf{v}$ and the other half to an optimal value function $\mathbf{v}^{\prime} \neq \mathbf{v}$. Since a neural network trained on this dataset is likely to converge to $\frac{1}{2}\left(\mathbf{v}+\mathbf{v}^{\prime}\right)$, it would be advantageous if this value function was likewise optimal. Since we so far failed to either prove or disprove this question, we pose it as an open problem for future work:

Problem A.2. Is the set of (counterfactually) optimal value functions convex?
We were able to prove a related result regarding the corresponding set of optimal extensions:

Proposition A.3. The set $\mathrm{OE}\left(\sigma^{\mathcal{T}}\right.$, reach.) of reachably-optimal extensions is always convex. The same might not hold for $\operatorname{OE}\left(\sigma^{\mathcal{T}}, \mathrm{cf}.\right)$ and $\mathrm{OE}\left(\sigma^{\mathcal{T}}\right.$, univ. $)$.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-53.jpg?height=505&width=434&top_left_y=1078&top_left_x=840)

Figure A.6: A game with non-convex set of optimal extensions of $\sigma^{\mathcal{T}}$. Circles denote chance nodes with uniform strategy. The "direction" of the triangles denotes whether the node belongs to a maximizer (tip up) or minimizer (tip down).

Proof. The reachably-optimal part of the proposition follows from the fact that reachably-optimal extensions of $\sigma^{\mathcal{T}}$ are precisely the Nash equilibria of $G\left(\sigma^{\mathcal{T}}\right)$ (see the proof of Proposition 4.16) and the set of Nash equilibria is always convex.

We now construct a counterexample for the negative part of the proposition. Let $G$ be the game from Figure A.6, where all chance strategies are uniform, $\mathcal{T}$ consists of the topmost three nodes $r, c_{1}$, and $c_{2}$, and strategy $\sigma^{\mathcal{T}}$ of player one (the maximizer) at $r$ is to always play left.

Firstly, note that independently of how we extend $\sigma^{\mathcal{T}}$ into a $\sigma$ in the whole game (i.e., by defining $\sigma_{1}(I)$ and $\sigma_{2}(J)$ ), we will necessarily have $V_{2}^{\sigma}(J)=140$

[^21]Indeed, this is because $V_{2}^{\sigma}(J)$ is the average of $v_{2}^{\sigma}\left(h^{\prime}\right)$ and $v_{2}^{\sigma}(h)$, weighted by the reach probabilities of $h$ and $h^{\prime}$, and the first of these probabilities is always zero. This is the lowest value (or, rather, the only value) that player 2 can get (given player one's strategy at $r$ ) It follows that to see whether $\sigma$ is an optimal extension, we only need to verify the maximality condition for player one at $I$ (i.e., verify that $V_{1}^{\sigma_{1}, \sigma_{2}}(I)=\max _{\left.\rho_{1}\right|_{I} \in \Sigma_{1}} V_{1}^{\left.\rho_{1}\right|_{I}, \sigma_{2}}(I)$ ). If the condition holds, the extension is universally optimal. If it doesn't hold, the extension will be reachably optimal (since $I$ isn't reachable) but not counterfactually optimal (since $I$ is counterfactually reachable).

First, note that no matter which strategy player two picks, we always have

$$
\max _{\left.\rho_{1}\right|_{I} \in \Sigma_{1}} V_{1}^{\rho_{1}, \sigma_{2}}(I) \geq V_{1}^{(0,1), \sigma_{2}}(I)=\frac{1}{2} \cdot 1+\frac{1}{2} \cdot 0=0.5
$$

Moreover, the optimal strategy at $I$ depends the strategy player two selects at $J$. If $\sigma_{2}(J)$ is s.t. $v_{1}^{\sigma_{2}}(h)=0$, any strategy at $I$ will be optimal (since all strategies will yield $\left.V_{1}^{\cdot, \sigma_{2}}(I)=0.5\right)$. If $\sigma_{2}(J)$ is s.t. $v_{1}^{\sigma_{2}}(h)<0$, the only optimal strategy at $I$ will be $(0,1)$ (otherwise we would get $\left.V_{1}^{\cdot, \sigma_{2}}(I)<0.5\right)$. Analogously, when $\sigma_{2}(J)$ is s.t. $v_{1}^{\sigma_{2}}(h)>0$, the only optimal strategy at $I$ will be $(1,0)$ (which yields $\left.V_{1}^{\cdot, \sigma_{2}}(I)>0.5\right)$.

We define $\sigma_{2}^{A}(J)=\left(\frac{1}{2}, \frac{1}{2}\right), \sigma_{1}^{A}(I)=\left(\frac{1}{2}, \frac{1}{2}\right)$ and $\sigma^{B}(J)=(0,1), \sigma^{A}(I)=(0,1)$. By the above observation, both $\sigma^{A}$ and $\sigma^{B}$ are universally-optimal extensions of $\sigma^{\mathcal{T}}$. However, consider the strategy $\sigma^{M}:=\frac{1}{2} \sigma^{A}+\frac{1}{2} \sigma^{B}$ (i.e., the $\frac{1}{2}$-convex combination of the two strategies). We have $\sigma_{1}^{M}(I)=\left(\frac{1}{2}, \frac{1}{2}\right)$ and $\sigma_{2}^{M}(J)=\left(\frac{1}{4}, \frac{3}{4}\right)$. By the above observation, this strategy is not counterfactually optimal.

However, the "slices" of the set of optimal extensions (for a fixed strategy of one of the players) are convex:

Proposition A.4. For $\sigma^{\mathcal{T}}$ and fixed $\sigma_{-p} \supset \sigma_{-p}^{\mathcal{T}}$, the set of extensions $\sigma_{p} \supset \sigma_{p}^{\mathcal{T}}$ for which $\sigma=\left(\sigma_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)$ is universally (resp. counterfactually) optimal is convex.

Proposition A. 4 provides some evidence that the answer to Problem A. 2 might be positive. And while the matter is by no means settled, this would correspond to the fact that we have, so far, not encountered any practical issues with neural networks generating value functions which would lead to exploitable strategies when used for depth-limited solving.

Proof. Let $\lambda \in[0,1]$ and suppose that the strategy profiles $\left(\rho_{p}, \sigma_{-p}\right)$ and $\left(\nu_{p}, \sigma_{-p}\right)$ are both universally-optimal extensions of $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$. Clearly, the convex combination $\mu_{p}:=\lambda \rho_{p}+(1-\lambda) \nu_{p}$ also extends $\sigma_{p}^{\mathcal{T}}$. To prove the proposition, it remains to show that the strategy profile $\left(\mu_{p}, \sigma_{-p}\right)$ is a universally-optimal extension of $\sigma^{\mathcal{T}}$.

Firstly, recall that the convex combination of strategies is defined in such a way that the mapping $\sigma_{p} \mapsto v_{i}^{\sigma_{p}, \sigma_{-p}}(h), i \in\{1,2\}$, is convex for any fixed $\sigma_{-p}$ and $h$. It follows that for any $\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{i}$, the functions $\sigma_{p} \mapsto V_{i}^{\sigma}(I)=$ $\sum_{h \in I} P^{\sigma}(h \mid I) v_{i}^{\sigma_{p}, \sigma_{-p}}(h)$ are likewise convex for any (since the beliefs $P^{\sigma}(h \mid I)$ only depend on $\left.\sigma^{\mathcal{T}}\right)$. By the assumptions on $\rho_{p}$ and $\nu_{p}$, we have $V_{i}^{\rho_{p}, \sigma_{-p}}(I)=$
$V_{i}^{\nu_{p}, \sigma_{-p}}(I)=\max _{\sigma_{p} \in \Sigma_{p}} V_{i}^{\sigma_{p}, \sigma_{-p}}(I)$ for all $\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{i}, i=1,2$. This concludes the proof, since the convexity of $V_{i}^{(\cdot), \sigma_{-p}}(I)$ implies that $\mu_{p}$ satisfies the optimality condition for all $I \subset \mathcal{Z}^{\mathcal{T}}$. (The proof in the counterfactually-optimal case is analogous.)

## B. Proofs

In this section, we give the proofs of the theoretical results presented in the main text.

Lemma 3.2. Let $(T, \sqsubset)$ be a finite tree, $Z \subset T$ its leaves, $f: Z \rightarrow \mathbb{R}$, and $P: T^{2} \rightarrow[0,1]$ a function s.t. $P(t, t)=1, P(t, s)=0$ when $\neg(t \sqsubset s)$, and $P(s, u)=P(s, t) P(t, u)$ when $s \sqsubset t \sqsubset u$. Then the following are equivalent for $F: T \rightarrow \mathbb{R}$ :
(a) $F(t)=\sum_{z \in Z} P(t, z) f(z)$ for $t \in T$,
(b) $F(z)=f(z)$ on $Z$ and $F(s)=\sum_{t \in \operatorname{ims}(s)} P(s, t) F(t)$ for $s \in T \backslash Z$ (where $\operatorname{ims}(s)$ is the set of all immediate successors of $s$ in $T)$,
(c) $F(z)=f(z)$ on $Z$ and $F(s)=\sum_{t \in L} P(s, t) F(t)$ for $s \in T \backslash Z$, whenever $L$ is a slice through $T$ below $s$.

Proof. $(c) \equiv(b)$ : Note that in $(c)$, the sum could equivalently be over $L^{\prime}:=\{t \in$ $L \mid s \sqsubset t\}$ (as $P(s, t)=0$ for the remaining $t \in L$ ). (c) trivially implies (b) by choosing $L:=$ any slice containing $\operatorname{ims}(() s)$. In the opposite direction, let $L$ be a slice below $s, L^{\prime}$ as above, and denote $L_{0}:=\{s\}, L_{i+1}:=\left(L \cap L_{i}\right) \cup \bigcup\{\operatorname{ims}(t) \mid$ $\left.t \in L_{i} \backslash L\right\}$. For $k:=\max _{t \in L}$ length $(t)$ - length $(s)$, we have $L_{k}=L^{\prime}$. If $F$ satisfies (b), it (trivially) satisfies $F(s)=\sum_{t \in L_{0}} P(s, t) F(t)$. Suppose that ( $c$ ) holds for $L_{i}$. Using the assumption from (b) on any $t \in L_{i} \backslash L^{\prime}$, we get that (c) holds for $L_{i+1}$. This shows that $(c)$ holds for $L_{k}=L^{\prime}$, thus proving (c).
$(a) \Longrightarrow(b)$ : Since $P(z, z)=1, F$ from $(a)$ coincides with $f$ on $Z$. Moreover, for any $s \in T \backslash Z$, we have

$$
\begin{align*}
F(s) & =\sum_{z \in Z} P(s, z) f(z)=\sum_{t \in \operatorname{ims}(s)} \sum_{t \sqsubset z \in Z} P(s, z) f(z)  \tag{B.1}\\
& =\sum_{t \in \operatorname{ims}(s)} \sum_{t \sqsubset z \in Z} P(s, t) P(t, z) f(z)  \tag{B.2}\\
& =\sum_{t \in \operatorname{ims}(s)} P(s, t) \sum_{t \sqsubset z \in Z} P(t, z) f(z)=\sum_{t \in \operatorname{ims}(s)} P(t), \tag{B.3}
\end{align*}
$$

which proves (b).
$(b) \Longrightarrow(a)$ : When $F$ coincides with $f, F$ satisfies (a) (i.e., $F(t)=$ $\left.\sum_{z \in Z} P(t, z) f(z)\right)$ for $t \in Z$. When $F$ satisfies $(a)$ for all $t \in \operatorname{ims}(() s)$, we have

$$
\begin{align*}
F(s) & =\sum_{t \in \operatorname{ims}(s)} P(s, t) F(t)=\sum_{t \in \operatorname{ims}(s)} P(s, t) \sum_{z \in Z} P(t, z) f(z)  \tag{B.4}\\
& =\sum_{t \in \operatorname{ims}(s)} P(s, t) \sum_{t \sqsubset z \in Z} P(t, z) f(z)  \tag{B.5}\\
& =\sum_{t \in \operatorname{ims}(s)} \sum_{t \sqsubset z \in Z} P(s, t) P(t, z) f(z)  \tag{B.6}\\
& =\sum_{t \in \operatorname{ims}(s)} \sum_{t \sqsubset z \in Z} P(s, z) f(z)=\sum_{z \in Z} P(s, z) f(z) . \tag{B.7}
\end{align*}
$$

Thus, by backwards induction, $F$ satisfies (a) for all $t \in T$.
Lemma 3.3 (Factorization of Infoset Reach Probabilities). For any $\sigma \in \Sigma$, $p \in \mathcal{N}$, and $I \in \mathcal{I}_{p}$, we have

$$
\begin{align*}
P_{p}^{\sigma}(I) & =P_{p}^{\sigma}(h) \text { for each } h \in I,  \tag{3.20}\\
P_{-p}^{\sigma}(I) & =\sum_{h \in I} P_{-p}^{\sigma}(h), \text { and }  \tag{3.21}\\
P^{\sigma}(I) & =P_{p}^{\sigma}(I) P_{-p}^{\sigma}(I) \tag{3.22}
\end{align*}
$$

Proof. Let $I \in \mathcal{I}_{p}$. Since $I$ is assumed to be thin, it is equal to its upper frontier and we have $P^{\mu}(I)=\sum_{h \in I} P^{\mu}(h)$ for any $\mu \in \Sigma$. By perfect recall, the sequence of actions taken by $p$ on the way to $h$ is the same for every $h \in I$, and likewise with the sequence of infosets encountered. As a result, the number $P_{p}^{\mu}(h)$ is the same for every $h \in I$. Together with the product formula (3.7), this gives the first part of the lemma:

$$
\begin{align*}
P_{-p}^{\sigma}(I) & =\max _{\rho_{p} \in \Sigma_{p}} P^{\rho_{p}, \sigma_{-p}}(I)=\max _{\rho_{p} \in \Sigma_{p}} \sum_{h \in I} P^{\rho_{p}, \sigma_{-p}}(h)  \tag{B.8}\\
& =\max _{\rho_{p} \in \Sigma_{p}} \sum_{h \in I} P_{p}^{\rho_{p}, \sigma_{-p}}(h) P_{-p}^{\rho_{p}, \sigma_{-p}}(h)  \tag{B.9}\\
& =\max _{\rho_{p} \in \Sigma_{p}} P_{p}^{\rho_{p}, \sigma_{-p}}\left(h_{0}\right) \sum_{h \in I} P_{-p}^{\sigma}(h)=\sum_{h \in I} P_{-p}^{\sigma}(h) \tag{B.10}
\end{align*}
$$

(where $h_{0}$ is an arbitrary element of $I$ ). The last equation holds because the maximum is realized by any strategy where $p$ always selects, on the path to $I$, the action leading towards $I$ (which leads to $P_{p}^{\rho_{p}, \sigma_{-p}}(h)=1$ ).

For the second equality, we have

$$
\begin{equation*}
P_{p}^{\sigma}(I)=P_{p}^{\sigma}\left(h_{0}\right) \max _{\rho_{-p} \in \Sigma_{-p}} \max _{\rho_{c} \in \Sigma_{c}} \sum_{h \in I} P_{-p}^{\sigma_{p}, \rho_{-p}, \rho_{c}}(h) \geq P_{p}^{\sigma}\left(h_{0}\right) \tag{B.11}
\end{equation*}
$$

(witnessed by picking a strategy for chance and $-p$ that plays to reach $h_{0}$ ). In the opposite direction, $P_{p}^{\sigma}(I)$ cannot be strictly higher than $P_{p}^{\sigma}\left(h_{0}\right)$ because -
no matter what chance and player $-p$ do - there is a $1-P_{p}^{\sigma}\left(h_{0}\right)$ chance that player $p$ will take some action that leads away from $h_{0}$, and hence (by perfect recall) away from $I$.

Finally, we use first two equations to get Equation (3.22):

$$
\begin{align*}
P^{\sigma}(I) & =\sum_{h \in I} P^{\sigma}(h)=\sum_{h \in I} P_{p}^{\sigma}(h) P_{-p}^{\sigma}(h)=\sum_{h \in I} P_{p}^{\sigma}(I) P_{-p}^{\sigma}(h)  \tag{B.12}\\
& =P_{p}^{\sigma}(I) \sum_{h \in I} P_{-p}^{\sigma}(h)=P_{p}^{\sigma}(I) P_{-p}^{\sigma}(I) \tag{B.13}
\end{align*}
$$

Lemma 3.5 (Equivalent definitions of the infoset belief). Let $\sigma \in \Sigma, I \in \mathcal{I}_{p}$. (1) The limit defining $P^{\sigma}(h \mid I)$ always exists. (2) For cf. reachable $I, P^{\sigma}(h \mid I)=$ $\lim _{n} P^{\left(1-\frac{1}{n}\right) \sigma+\frac{1}{n} \text { unif }}(h \mid I)=\frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}$. (3) For reachable $I$,

$$
\begin{equation*}
P^{\sigma}(h \mid I)=\lim _{n \rightarrow \infty} P^{\left(1-\frac{1}{n}\right) \sigma+\frac{1}{n} \text { unif }}(h \mid I)=\frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}=\frac{P^{\sigma}(h)}{P^{\sigma}(I)} \tag{3.24}
\end{equation*}
$$

Proof. (1): For $x \in(0,1]$, denote $\sigma^{x}:=(1-x) \sigma+x$ unif. Our goal is to show that the limit $P^{\sigma^{x}}(h) / \sum_{g \in I} P^{\sigma^{x}}(I)$ exists. For every $g \in I, P^{\sigma^{x}}(g)$ is the product of $\sigma^{x}(k, a)=(1-x) \sigma(k, a)+x /|\mathcal{A}(k)|, k a \sqsubset g$. If $n=n_{g}$ of the $\sigma(k, a)$-s is non-zero, the product will be equal to $C x^{n}(1+x \varphi(x))$, where $C=C_{g}>0$ is some constant (the product of non-zero $\sigma(k, a)$-s and $1 /|\mathcal{A}(k)|$-s for those $k$ where $\sigma(k, a)=0$ ) and $\varphi=\varphi_{g}$ is some polynomial of $x$. By denoting $m:=\min _{g \in I} n_{g}$, we can further rewrite the product as $x^{m} x \psi_{g}(x)$ whenever $n_{g}>m$. We then have

$$
\begin{equation*}
\frac{P^{\sigma^{x}}(h)}{\sum_{g \in I} P^{\sigma^{x}}(h)}=\frac{x^{m}}{x^{m}} \cdot \frac{C_{h} x^{n_{h}-m}\left(1+x \varphi_{g}(x)\right)}{\sum_{g, n_{g}=m} C_{g}\left(1+x \varphi_{g}(x)\right)+x \sum_{g, n_{g}>m} \psi_{g}(x)} \tag{B.14}
\end{equation*}
$$

As $x$ tends to 0 , the limit of this expression is either 0 (when $n_{h}>m$ ) or $\frac{C_{h}}{\sum_{g, n_{g}>m} C_{g}}$ (when $n_{h}=m$ ). This proves (1).
(2): Recall that for $I \in \mathcal{I}_{p}, P_{p}^{\sigma}(g)=: P_{p}^{\sigma}(I)$ is the same for all $g \in I$. We thus have

$$
\begin{equation*}
\frac{P^{\sigma^{x}}(h)}{P^{\sigma^{x}}(I)}=\frac{P_{p}^{\sigma^{x}}(h) P_{-p}^{\sigma^{x}}(h)}{P_{p}^{\sigma^{x}}(I) P_{-p}^{\sigma^{x}}(I)}=\frac{P_{-p}^{\sigma^{x}}(h)}{P_{-p}^{\sigma^{x}}(I)} . \tag{B.15}
\end{equation*}
$$

As $x \rightarrow 0$, we have $P_{p}^{\sigma^{x}}(h) \longrightarrow P_{-p}^{\sigma}(h)$ and $P_{p}^{\sigma^{x}}(I) \longrightarrow P_{-p}^{\sigma}(I)$. Since $I$ being counterfactually reachable means that $P_{-p}^{\sigma^{x}}(I)>0$, we get $\frac{P_{-p}^{\sigma^{x}}(h)}{P_{-p}^{\sigma}(I)} \longrightarrow \frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}$ as $x \rightarrow 0$.
(3): When $P^{\sigma}(I)>0$, we even get $\frac{P^{\sigma^{x}}(h)}{P^{\sigma^{x}}(I)} \longrightarrow \frac{P^{\sigma}(h)}{P^{\sigma}(I)}$.

Lemma 3.6 (Properties of Infoset Reach Probabilities). For any $\sigma \in \Sigma$ and $I \sqsubset J \sqsubset K$ in $\mathcal{I}_{p}$, we have:

$$
\begin{equation*}
P^{\sigma}(I, J)=\lim _{n \rightarrow \infty} \frac{P^{\sigma^{n}}(J)}{P^{\sigma^{n}}(I)}, P_{p}^{\sigma}(I, J)=\lim _{n \rightarrow \infty} \frac{P_{p}^{\sigma^{n}}(J)}{P_{p}^{\sigma^{n}}(I)}, \tag{1}
\end{equation*}
$$

$$
\text { and } P_{-p}^{\sigma}(I, J)=\lim _{n \rightarrow \infty} \frac{P_{-p}^{\sigma^{n}}(J)}{P_{-p}^{\sigma^{n}(I)}} \text { (where } \sigma^{n} \text { denotes } \frac{n-1}{n} \sigma+\frac{1}{n} \text { unif). }
$$

(2) $P^{\sigma}(I, J)=P_{p}^{\sigma}(I, J) P_{-p}^{\sigma}(I, J)$.
(3) $P^{\sigma}(I, K)=P^{\sigma}(I, J) P^{\sigma}(J, K)$.

Proof. Whenever we manipulate limits in this proof, we should verify that this manipulation is correct and all the limits exist. To avoid repeating the calculations done in Lemma 3.5, we skip this part of the proof.
(1): For $n \in \mathbb{N}$, we have $P^{\sigma}(I, J)=$

$$
\begin{aligned}
& =\sum_{g \in I} P^{\sigma}(g \mid I) \sum_{g \sqsubset \in J} P^{\sigma}(g, h)=\sum_{g \in I} \lim _{n \rightarrow \infty} \frac{P^{\sigma^{n}}(g)}{P^{\sigma^{n}}(I)} \sum_{g \sqsubset h \in J} P^{\sigma}(g, h) \\
& =\lim _{n \rightarrow \infty} \sum_{g \in I} \frac{P^{\sigma^{n}}(g)}{P^{\sigma^{n}}(I)} \sum_{g \sqsubset h \in J} P^{\sigma}(g, h) \\
& =\lim _{n \rightarrow \infty} \sum_{g \in I} \frac{P^{\sigma^{n}}(g)}{P^{\sigma^{n}}(I)}\left(\sum_{g \sqsubset h \in J}\left(P^{\sigma^{n}}(g, h)+P^{\sigma}(g, h)-P^{\sigma^{n}}(g, h)\right)\right) \\
& =\lim _{n \rightarrow \infty}\left(\sum_{g \in I} \sum_{g \sqsubset h \in J} \frac{P^{\sigma^{n}}(g)}{P^{\sigma^{n}}(I)} P^{\sigma^{n}}(g, h)\right)+\sum_{g \in I}\left[\left(\lim _{n \rightarrow \infty} \frac{P^{\sigma^{n}}(g)}{P^{\sigma^{n}}(I)}\right)\left(\sum_{g \sqsubset h \in J} 0\right)\right] \\
& =\lim _{n \rightarrow \infty} \frac{P^{\sigma^{n}}(I)}{P^{\sigma^{n}}(I)} .
\end{aligned}
$$

The proofs of the other two identities are analogous.
(2): By (1), we have

$$
\begin{aligned}
& P^{\sigma}(I, J) \stackrel{n \rightarrow \infty}{\longleftrightarrow} \\
& \frac{P^{\sigma^{n}}(J)}{P^{\sigma^{n}}(I)}=\frac{P_{p}^{\sigma^{n}}(J) P_{-p}^{\sigma^{n}}(J)}{P_{p}^{\sigma^{n}}(I) P_{-p}^{\sigma^{n}}(I)}= \\
&=\frac{P_{p}^{\sigma^{n}}(J)}{P_{p}^{\sigma^{n}}(I)} \frac{P_{-p}^{\sigma^{n}}(J)}{P_{-p}^{\sigma^{n}}(I)} \xrightarrow{n \rightarrow \infty} P_{p}^{\sigma}(I, J) P_{-p}^{\sigma}(I, J) .
\end{aligned}
$$

(3): By (1), we have

$$
P^{\sigma}(I, K) \stackrel{n \rightarrow \infty}{\longleftrightarrow} \frac{P^{\sigma^{n}}(K)}{P^{\sigma^{n}}(I)}=\frac{P^{\sigma^{n}}(K) P^{\sigma^{n}}(J)}{P^{\sigma^{n}}(J) P^{\sigma^{n}}(I)} \xrightarrow{n \rightarrow \infty} P^{\sigma}(I, J) P^{\sigma}(J, K) .
$$

Lemma 3.8. For any $\sigma \in \Sigma, I \in \mathcal{I}_{p}$ s.t. $\mathcal{P}(I)=p$, and $a \in \mathcal{A}(I)$, we have

$$
\begin{equation*}
Q^{\sigma}(I, a)=\sum_{J \in \operatorname{ims}(I, a)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J) \tag{3.35}
\end{equation*}
$$

Proof. Let $\sigma, I$, and $a$ be as in the statement of the lemma. We will expand the right-hand side of the equation to get $Q^{\sigma}(I, a)$. For any $J \in \operatorname{ims}(I, a)$, denote by $J^{\prime}$ the subset of $I$ for which $\left\{h a \mid h \in J^{\prime}\right\}=J$. First, for any $g \in I$, there is either no $h \in J$ s.t. $g \sqsubset h$ (when $g \notin J^{\prime}$ ) or exactly one such $h(h=g a$, when $\left.g \in J^{\prime}\right)$, in which case $P_{-p}^{\sigma}(g, h)=1$. As a result, the definition of $P_{-p}^{\sigma}(I, J)$ gives

$$
\begin{equation*}
P_{-p}^{\sigma}(I, J)=\sum_{g \in I} P^{\sigma}(g \mid I) \sum_{g \sqsubset h \in J} P_{-p}^{\sigma}(g, h)=\sum_{g \in J^{\prime}} P^{\sigma}(g \mid I)=P^{\sigma}\left(J^{\prime} \mid I\right) . \tag{B.16}
\end{equation*}
$$

To expand $V^{\sigma}(J)$, note first that for any $g a \in J$, we have $P^{\sigma}(g a \mid J)=P^{\sigma}\left(g \mid J^{\prime}\right)$. (Indeed, this is immediate when $J$ is reachable, since $\sigma_{p}(g, a)$ is the same for all $g \in I$. For non-reachable $J$, we just take the limit over the amount of uniformly random noise added to $\sigma$.) We then have

$$
\begin{equation*}
V^{\sigma}(J)=\sum_{g a \in J} P^{\sigma}(g a \mid J) v_{p}^{\sigma}(g a)=\sum_{g \in J^{\prime}} P^{\sigma}\left(g \mid J^{\prime}\right) q_{p}^{\sigma}(g, a) \tag{B.17}
\end{equation*}
$$

Putting the two equations together, we get

$$
\begin{align*}
\sum_{J \in \operatorname{ims}(I, a)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J) & =\sum_{J^{\prime} \subset I} P^{\sigma}\left(J^{\prime} \mid I\right) \sum_{g \in J^{\prime}} P^{\sigma}\left(g \mid J^{\prime}\right) q_{p}^{\sigma}(g, a)  \tag{B.18}\\
& =\sum_{J^{\prime} \subset I} \sum_{g \in J^{\prime}} P^{\sigma}\left(J^{\prime} \mid I\right) P^{\sigma}\left(g \mid J^{\prime}\right) q_{p}^{\sigma}(g, a)  \tag{B.19}\\
& =\sum_{J^{\prime} \subset I} \sum_{g \in J^{\prime}} P^{\sigma}(g \mid I) q_{p}^{\sigma}(g, a)  \tag{B.20}\\
& =\sum_{g \in I} P^{\sigma}(g \mid I) q_{p}^{\sigma}(g, a)=Q^{\sigma}(I, a) \tag{B.21}
\end{align*}
$$

Theorem 1 (Characterization of $V^{\sigma}$ ). Suppose that terminal infosets are always singleton. Then for any $p \in \mathcal{N}$ and $\sigma \in \Sigma$, we have $V_{p}^{\sigma}$ (root $)=u_{p}(\sigma)$ and the functions $V^{\sigma}: \mathcal{I}_{p} \rightarrow \mathbb{R}$ and $Q^{\sigma}(I, a):=\sum_{\operatorname{ims}(I, a)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J)$ have the following properties:
(1) $V^{\sigma}(I)=\sum_{h \in I} P^{\sigma}(h \mid I) v_{p}^{\sigma}(h)$.
(2) $V^{\sigma}(I)=\sum_{z \in \mathcal{Z}} P^{\sigma}(I,\{z\}) u_{p}(z)$.
(3) $V_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and for non-terminal $I \in \mathcal{I}_{p}$, we have

$$
\begin{equation*}
V^{\sigma}(I)=\sum_{J \in \operatorname{ims}(I)} P^{\sigma}(I, J) V^{\sigma}(J) \tag{3.36}
\end{equation*}
$$

(3') $V_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and for non-terminal $I$, we have

$$
V^{\sigma}(I)= \begin{cases}\sum_{a \in \mathcal{A}_{p}(I)} \sigma_{p}(I, a) Q^{\sigma}(I, a) & \text { when } p \text { acts in } I \\ \sum_{J \in \operatorname{ims}(I)} P_{-p}^{\sigma}(I, J) V^{\sigma}(J) & \text { when } p \text { doesn't act in } I\end{cases}
$$

(4) $V_{p}^{\sigma}=u_{p}$ on $\mathcal{Z}$ and for non-terminal $I \in \mathcal{I}_{p}$ and any slice $\mathcal{L}$ through $\mathcal{I}_{p}$ below $I$, we have $V^{\sigma}(I)=\sum_{J \in \mathcal{L}} P^{\sigma}(I, J) V^{\sigma}(J)$.

Moreover, each of these conditions can be used as an equivalent definition of $V_{p}^{\sigma}$ (i.e., it automatically implies all the other properties).

Proof. $V_{p}^{\sigma}($ root $)=u_{p}(\sigma)$ follows from (1) since $v_{p}^{\sigma}($ root $)=u_{p}(\sigma)$ by Lemma 3.1.
$(3) \Longleftrightarrow\left(3^{\prime}\right)$ holds because the two formulas are equivalent. Indeed, when $J \in \operatorname{ims}(I)$ and $\mathcal{P}(I) \neq p$, we have $P_{p}^{\sigma}(I, J)=1$, so $P^{\sigma}(I, J)=P_{-p}^{\sigma}(I, J)$. When $\mathcal{P}(p)=p$, Lemma 3.8 yields

$$
\begin{align*}
\sum_{J \in \operatorname{ims}(I)} P^{\sigma}(I, J) V^{\sigma}(J) & =\sum_{a \in \mathcal{A}(I)} \sum_{J \in \operatorname{ims}(I, a)} \sigma_{p}(I, a) P_{-p}^{\sigma}(I, J) V^{\sigma}(J)  \tag{B.22}\\
& =\sum_{a \in \mathcal{A}(I)} \sigma_{p}(I, a) Q^{\sigma}(I, a) \tag{B.23}
\end{align*}
$$

$(2) \Longleftrightarrow(3) \Longleftrightarrow(4)$ : These equivalences follow from Lemma 3.2. Indeed, we get this by applying the lemma to the tree $(T, \sqsubset):=\left(\mathcal{I}_{p}, \sqsubset\right), P:=P^{\sigma}(\cdot, \cdot)$ (which, by Lemma3.6, satisfies $P^{\sigma}(I, K)=P^{\sigma}(I, J) P^{\sigma}(J, K)$ whenever $I, J, K \in$ $\mathcal{I}_{p}$ are s.t. $\left.I \sqsubset J \sqsubset K\right), f(\{z\}):=u_{p}(z)$, and $F:=V^{\sigma}$.
$(1) \Longleftrightarrow(3)$ : Denote by $V^{\prime}$ the $V^{\sigma}$ from (1) and by $V^{\prime \prime}$ the $V^{\sigma}$ from (3). We will use backwards induction on $\mathcal{I}_{p}$ to show that $V^{\prime}=V^{\prime \prime}$. First, when $z \in \mathcal{Z}$, the infoset $I=\{z\}$ satisfies

$$
\begin{equation*}
V^{\prime}(I)=\sum_{h \in I} P^{\sigma}(h \mid I) v_{p}^{\sigma}(h)=P^{\sigma}(z \mid\{z\}) v_{p}^{\sigma}(z)=1 \cdot u_{p}(z)=V^{\prime \prime}(I) \tag{B.24}
\end{equation*}
$$

In particular, the two functions are equal on the leaves of $\mathcal{I}_{p}$. Second, suppose that $V^{\prime}(J)=V^{\prime \prime}(J)$ for all $J \in \operatorname{ims}(I)$. We then have $V^{\prime \prime}(I)=$

$$
\begin{aligned}
& =\sum_{J \in \operatorname{ims}(I)} P^{\sigma}(I, J) V^{\prime \prime}(J)=\sum_{J \in \operatorname{ims}(I)} P^{\sigma}(I, J) V^{\prime}(J) \\
& =\sum_{J \in \operatorname{ims}(I)} P^{\sigma}(I, J) \sum_{h \in J} P^{\sigma}(h \mid J) v_{p}^{\sigma}(h)=\sum_{J \in \operatorname{ims}(I)} \sum_{h \in J} P^{\sigma}(I, J) P^{\sigma}(h \mid J) v_{p}^{\sigma}(h) \\
& \stackrel{(X)}{=} \sum_{g \in I} \sum_{h \in \operatorname{ims}(g)} P^{\sigma}(g \mid I) P^{\sigma}(g, h) v_{p}^{\sigma}(h)=\sum_{g \in I} P^{\sigma}(g \mid I) \sum_{h \in \operatorname{ims}(g)} P^{\sigma}(g, h) v_{p}^{\sigma}(h) \\
& \stackrel{(Y)}{=} \sum_{g \in I} P^{\sigma}(g \mid I) v_{p}^{\sigma}(g)=V^{\prime}(I) .
\end{aligned}
$$

(X) follows from the fact that $P^{\sigma}(I, J) P^{\sigma}(h \mid I)=P^{\sigma}(g \mid I) P^{\sigma}(g, h)$ holds for the unique $g \in I$ for which $g \sqsubset h$. (This is easy for fully-mixed $\sigma$ since $P^{\sigma}(I, J) P^{\sigma}(h \mid I)=\frac{P^{\sigma}(J)}{P^{\sigma}(I)} \frac{P^{\sigma}(h)}{P^{\sigma}(I)}=\frac{P^{\sigma}(h)}{I}$ and $P^{\sigma}(g \mid I) P^{\sigma}(g, h)=\frac{P^{\sigma}(g)}{P^{\sigma}(I)} P^{\sigma}(g, h)=$ $\frac{P^{\sigma}(h)}{P^{\sigma}(I)}$. For general $\sigma$, we take the limit over $\sigma^{n}=\frac{n-1}{n} \sigma+\frac{1}{n}$ unif.) For (Y), we used Lemma 3.1.

Theorem 2 (Properties of $V_{\mathrm{cf}}^{\sigma}$ ). For any $\sigma \in \Sigma, p \in \mathcal{N}$, and non-terminal $I \in \mathcal{I}_{p}$, we have:
(1) $V_{p, \text { cf }}^{\sigma}($ root $)=u_{p}(\sigma)$.
(2) $V_{\mathrm{cf}}^{\sigma}(I)=P_{-p}^{\sigma}(I) V^{\sigma}(I)$.
(3) For any slice through $\mathcal{I}_{p}$ below $I$, we have $V_{\mathrm{cf}}^{\sigma}(I)=\sum_{J \in \mathcal{L}} P_{p}^{\sigma}(I, J) V_{\mathrm{cf}}^{\sigma}(J)$.
(4) (a) For terminal $Z \in \mathcal{I}_{p}, V_{\mathrm{cf}}^{\sigma}(Z)=\sum_{z \in Z} P_{-p}^{\sigma}(z) u_{p}(z)$.
(b) When $\mathcal{P}(I)=p$, we have

$$
V_{\mathrm{cf}}^{\sigma}(I)=\sum_{a \in \mathcal{A}(I)} \sigma_{p}(I, a) Q_{\mathrm{cf}}^{\sigma}(I, a)=\sum_{a \in \mathcal{A}(I)} \sigma_{p}(I, a) \sum_{J \in \operatorname{ims}(I, a)} V_{\mathrm{cf}}^{\sigma}(J)
$$

(c) When $\mathcal{P}(I) \neq p$, we have $V_{\mathrm{cf}}^{\sigma}(I)=\sum_{J \in \operatorname{ims}(I)} V_{\mathrm{cf}}^{\sigma}(J)$.

Proof. (1): This holds because $V_{p, \text { cf }}^{\sigma}($ root $)=P_{-p}^{\sigma}($ root $) v_{p}^{\sigma}$ (root) $=1 \cdot v_{p}^{\sigma}($ root $)$ (and Lemma 3.1 gives $v_{p}^{\sigma}($ root $\left.)=u_{p}(\sigma)\right)$.
(2): First, note that for any $I$ and $\sigma, P_{-p}^{\sigma}(I) P^{\sigma}(h \mid I)=P_{-p}^{\sigma}(I)$. (For counterfactually reachable $I$, this holds because $P^{\sigma}(h \mid I)=\frac{P_{-p}^{\sigma}(h)}{P_{-p}^{\sigma}(I)}$ by Lemma 3.5. For general $I$, the result holds because $P^{\sigma}(h \mid I)=\lim _{n} P^{\sigma^{n}}(h \mid I)$, where $\sigma^{n}$ is fully mixed and $\sigma^{n} \rightarrow \sigma$.) Consequently, we have

$$
P_{-p}^{\sigma}(I) V^{\sigma}(I)=P_{-p}^{\sigma}(I) \sum_{h \in I} P^{\sigma}(h \mid I) v_{p}^{\sigma}(h)=\sum_{h \in I} P_{-p}^{\sigma}(h) v_{p}^{\sigma}(h) .
$$

(3): First, observe that

$$
P_{-p}^{\sigma}(I) P^{\sigma}(I, J)=P_{-p}^{\sigma}(I) P_{-p}^{\sigma}(I, J) P_{p}^{\sigma}(I, J)=P_{-p}^{\sigma}(J) P_{p}^{\sigma}(I, J) .
$$

Combining this fact with Theorem 1, we get

$$
\begin{aligned}
V_{\mathrm{cf}}^{\sigma}(I) & =P_{-p}^{\sigma}(I) V^{\sigma}(I)=\sum_{J \in \mathcal{L}} P_{-p}^{\sigma}(I) P^{\sigma}(I, J) V^{\sigma}(J) \\
& =\sum_{J \in \mathcal{L}} P_{-p}^{\sigma}(J) P_{p}^{\sigma}(I, J) V^{\sigma}(J)=\sum_{J \in \mathcal{L}} P_{p}^{\sigma}(I, J) V_{\mathrm{cf}}^{\sigma}(J) .
\end{aligned}
$$

(4): (a) holds because $V_{\mathrm{cf}}^{\sigma}(Z)=\sum_{z \in Z} P_{-p}^{\sigma}(z) v_{p}^{\sigma}(z)$ (by definition of $V_{\mathrm{cf}}^{\sigma}$ ) and $v_{p}^{\sigma}(z)=u_{p}(z)$ whenever $z \in \mathcal{Z}$. To prove (b), observe first that

$$
\begin{aligned}
Q_{\mathrm{cf}}^{\sigma}(I, a) & =\sum_{h \in \mathcal{A}(I)} q_{p, \mathrm{cf}}^{\sigma}(h, a)=\sum_{h \in \mathcal{A}(I)} v_{p, \mathrm{cf}}^{\sigma}(h a) \\
& =\sum_{J \in \mathrm{ims}(I, a)} \sum_{h a \in J} v_{p, \mathrm{cf}}^{\sigma}(h a)=\sum_{J \in \operatorname{ims}(I, a)} V_{p \mathrm{cf}}^{\sigma}(J) .
\end{aligned}
$$

Moreover, when $\mathcal{P}(I)=p$, we have $P_{p}^{\sigma}(I, J)=\sigma_{p}(I, a)$ for any $J \in \operatorname{ims}(I, a)$. Combining this with (3) applied to $\mathcal{L}:=$ 'some slice through $\mathcal{I}_{p}$ containing
$\operatorname{ims}(I)^{\prime}$, we get

$$
\begin{aligned}
V_{\mathrm{cf}}^{\sigma}(I) & \stackrel{(3)}{=} \sum_{a \in \mathcal{A}(I)} \sum_{J \in \operatorname{ims}(I, a)} P_{p}^{\sigma}(I, J) V_{\mathrm{cf}}^{\sigma}(J) \\
& =\sum_{a \in \mathcal{A}(I)} \sum_{J \in \operatorname{ims}(I, a)} \sigma_{p}(I, a) V_{\mathrm{cf}}^{\sigma}(J)=\sum_{a \in \mathcal{A}(I)} \sigma_{p}(I, a) Q_{\mathrm{cf}}^{\sigma}(I, a)
\end{aligned}
$$

The proof of (c) is analogous to the proof of (b), except that we use the fact that when $\mathcal{P}(I) \neq p$, we have $P_{p}^{\sigma}(I, J)=1$ for any $J \in \operatorname{ims}(I, a)$.

Lemma 4.10. Suppose that a value function satisfies $u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$ for every $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$. Then $\mathbf{v}$ preserves Nash equilibria of $G$.

Proof. Let $\mathbf{v}$ be a value function for $\mathcal{T}$ that satisfies the assumption of the lemma.

To prove that any $\left.\sigma^{\mathcal{T}} \in \operatorname{NE}(G)\right|_{\mathcal{T}}$ is a solution of $(\mathcal{T}, \mathbf{v})$, suppose that $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ admits some extension $\sigma^{*} \in \operatorname{NE}(G)$. Since each player can extend $\sigma^{\mathcal{T}}$ into a non-exploitable strategy in $G$, we have $\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right) \geq \operatorname{gv}_{p}(G)$ for both $p$, which implies that $\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)=\operatorname{gv}_{p}(G)$. If $p$ were to deviate from $\sigma^{\mathcal{T}}$ in $(\mathcal{T}, \mathbf{v})$ and use some strategy $\rho_{p}^{\mathcal{T}}$ instead, the opponent could use $\sigma_{-p}^{\mathcal{T}}$ which to ensure that $\operatorname{gv}_{-p}\left(G\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)\right) \geq \operatorname{gv}_{-p}(G)$. Since $G$ is zero-sum, we have $\operatorname{gv}_{p}\left(G\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)\right) \leq \operatorname{gv}_{p}(G)$. Taken together, this implies that

$$
u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\operatorname{gv}_{p}(G) \geq \operatorname{gv}_{p}\left(G\left(\rho_{1}^{\mathcal{T}}, \sigma_{2}^{\mathcal{T}}\right)\right)=u_{p}^{\mathbf{v}}\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)
$$

In other words, no player can gain utility in $(\mathcal{T}, \mathbf{v})$ by deviating from $\sigma^{\mathcal{T}}$, so $\sigma^{\mathcal{T}}$ is a solution of $(\mathcal{T}, \mathbf{v})$.

In the opposite direction, suppose that $\sigma^{\mathcal{T}}$ is a solution of $(\mathcal{T}, \mathbf{v})$. We will show that $\sigma^{\mathcal{T}}$ 's trunk exploitability is 0 , implying that it can be extended into a Nash equilibrium in $G$ (Proposition 4.5). Since $\sigma^{\mathcal{T}}$ is a solution of $(\mathcal{T}, \mathbf{v})$, it must satisfy

$$
\begin{equation*}
(\forall p)\left(\forall \rho_{p}^{\mathcal{T}} \in \Sigma_{p}^{\mathcal{T}}\right): u_{p}^{\mathbf{v}}\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right) \leq u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right) \tag{B.25}
\end{equation*}
$$

Using the assumption of the lemma, we get $\operatorname{gv}_{p}\left(G\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right) \leq \operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)\right.$ for each $\rho_{p}^{\mathcal{T}}$. Substituting $\rho_{p}^{\mathcal{T}}:=\sigma_{1}^{*} \mid \mathcal{T}$ for some $\sigma^{*} \in \operatorname{NE}(G)$, we get

$$
\begin{equation*}
\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right) \geq \operatorname{gv}_{p}\left(\left(G\left(\sigma_{1}^{*} \mid \mathcal{T}, \sigma_{-p}^{\mathcal{T}}\right)\right) \geq \operatorname{gv}_{p}(G)\right. \tag{B.26}
\end{equation*}
$$

Since the game is zero-sum, it follows that

$$
\begin{equation*}
\operatorname{gv}_{p}(G)=\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right) \tag{B.27}
\end{equation*}
$$

By Proposition 4.5, we have $\operatorname{expl}_{p}\left(\sigma_{p}^{\mathcal{T}}\right)=\operatorname{gv}_{p}(G)-\operatorname{gv}_{p}\left(G\left(\sigma_{p}^{\mathcal{T}}\right)\right)=0$.

Proposition 4.16 (Computing reachably-optimal values). Suppose that for each $\sigma^{\mathcal{T}}$, there is some $\sigma \in \operatorname{NE}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$ s.t. $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)=v_{1}^{\sigma}(h)$ holds for all $h \in \mathcal{Z}^{\mathcal{T}}$. Then $\mathbf{v}$ is a reachably optimal value function.

Proof. Recall that $\mathbf{v}$ is reachably optimal if every trunk strategy $\sigma^{\mathcal{T}}$ admits a reachably-optimal extension for which $\mathbf{v}^{\sigma^{\top}}(h)=v_{1}^{\sigma}(h)$ (for all histories in reachable infosets in $\mathcal{Z}^{\mathcal{T}}$. By the assumption of the proposition, each $\mathbf{v}^{\sigma^{\mathcal{T}}}(\cdot)$ coincides with $v_{1}^{\sigma}(\cdot)$ for some $\sigma \in N E\left(G\left(\sigma^{\mathcal{T}}\right)\right)$. To prove the result, we thus need to show that for any $\sigma^{\mathcal{T}}$, the corresponding extension $\sigma \in N E\left(G\left(\sigma^{\mathcal{T}}\right)\right)$ is reachably optimal.

Let $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ and let $\sigma$ be the corresponding element of $\mathrm{NE}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$. Suppose that $\sigma$ isn't reachably-optimal, i.e., suppose that there was some reachable $I_{0} \subset \mathcal{Z}^{\mathcal{T}}$ for which $p$ could increase their value of $I_{0}$ by switching to some strategy $\rho_{p} \neq \sigma_{p}$. Recall that as far as $p$ 's strategy is concerned, the values $V_{p}^{(\cdot), \sigma_{-p}}(I)$ only depend on what $p$ does below $I$ - not on what they do below other infosets in $\mathcal{Z}^{\mathcal{T}}$, nor on what they do in the trunk. Without loss of generality, we can therefore assume that (i) $V_{p}^{\rho_{p}, \sigma_{-p}}(I)=V_{p}^{\sigma_{p}, \sigma_{-p}}(I)$ for all $p$ 's infosets $I \neq I_{0}$ in $\mathcal{Z}^{\mathcal{T}}$ and (ii) $\rho_{p}$ is an extension of $\sigma_{p}^{\mathcal{T}}$. By Theorem 1 , we have

$$
\begin{aligned}
u_{p}\left(\rho_{p}, \sigma_{-p}\right) & =\sum_{I \in \mathcal{I}_{p}, I \subset \mathcal{Z} \mathcal{T}} P^{\rho_{p}, \sigma_{-p}}(I) V_{p}^{\rho_{p}, \sigma_{-p}}(I) \\
& =\sum_{I} P^{\sigma}(I) V_{p}^{\rho_{p}, \sigma_{-p}}(I) \\
& =P^{\sigma}\left(I_{0}\right) V_{p}^{\rho_{p}, \sigma_{-p}}\left(I_{0}\right)+\sum_{I \neq I_{0}} P^{\sigma}(I) V_{p}^{\sigma_{p}, \sigma_{-p}}(I) \\
& >P^{\sigma}\left(I_{0}\right) V_{p}^{\sigma_{p}, \sigma_{-p}}\left(I_{0}\right)+\sum_{I \neq I_{0}} P^{\sigma}(I) V_{p}^{\sigma_{p}, \sigma_{-p}}(I)=u_{p}\left(\sigma_{p}, \sigma_{-p}\right)
\end{aligned}
$$

(where the inequality required the assumption that the reach probability of $I_{0}$ is positive). This shows that $p$ could increase their utility in $G\left(\sigma^{\mathcal{T}}\right)$ by deviating from $\sigma_{p}$, which contradicts our assumption of $\sigma$ being Nash equilibrium. This shows that $\sigma$ must be a reachably-optimal extension of $\sigma^{\mathcal{T}}$, which concludes the proof.

Proposition 4.17 (Enabling utility calculation). Any reachably-optimal value function satisfies

$$
\begin{equation*}
u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\sum_{h \in \mathcal{Z} \mathcal{T}} P^{\sigma^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h)=\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right) \tag{4.6}
\end{equation*}
$$

Proof. Let $\sigma^{\mathcal{T}}$ be a trunk strategy and denote $\Sigma_{p}^{\downarrow}:=\left.\Sigma_{p}\right|_{\mathcal{H} \backslash \mathcal{T}}$. First, observe that since $G\left(\sigma^{\mathcal{T}}\right)$ is a zero-sum EFG, we have

$$
\begin{equation*}
\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)=\max _{\sigma_{1}^{\prime} \supset \sigma_{1}^{\mathcal{T}}} \min _{\sigma_{2}^{\prime} \supset \sigma_{2}^{\mathcal{T}}} u_{1}\left(\sigma^{\prime}\right)=\min _{\sigma_{2}^{\prime} \supset \sigma_{2}^{\mathcal{T}}} \max _{\sigma_{1}^{\prime} \supset \sigma_{1}^{\mathcal{T}}} u_{1}\left(\sigma^{\prime}\right) \tag{B.28}
\end{equation*}
$$

(by minimax theorem).

Second, denote $V^{\sigma^{\mathcal{T}}, \mathbf{v}}(I):=\sum_{h \in I} P^{\sigma^{\mathcal{T}}}(h \mid I) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h)$ for $I \subset \mathcal{Z}^{\mathcal{T}}$ and let $\sigma=\sigma^{\mathcal{T}} \cup \sigma^{\downarrow}$ be the reachably optimal extension of $\sigma^{\mathcal{T}}$ from the definition of $\mathbf{v}$ being reachably optimal. Expanding the definition of $u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)$ yields

$$
\begin{aligned}
u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right) & =\sum_{h \in \mathcal{Z}^{\mathcal{T}}} P^{\sigma^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h)=\sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\sigma^{\mathcal{T}}}(I) V^{\sigma^{\mathcal{T}}, \mathbf{v}}(I) \\
& =\sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\sigma}(I) V^{\sigma}(I)=u_{p}(\sigma)
\end{aligned}
$$

(The first identity holds by Theorem 1 . The third holds because an infoset $I$ is either reachable and $V_{p}^{\sigma^{\mathcal{T}}, \mathbf{v}}(I)=V_{p}^{\sigma}(I)$ or unreachable and $P^{\sigma^{\mathcal{T}}}(I) V^{\sigma^{\mathcal{T}}, \mathbf{v}}(I)=$ $P^{\sigma}(I) V^{\sigma}(I)$ because $P^{\sigma^{\top}}(I)=P^{\sigma}(I)=0$.)

Third, for $\left(\rho_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right) \in \Sigma_{1}^{\downarrow} \times \Sigma_{2}^{\downarrow}$, denote $f\left(\rho_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right):=u_{1}\left(\sigma_{1}^{\mathcal{T}} \cup \rho_{1}^{\downarrow}, \sigma_{2}^{\mathcal{T}} \cup \rho_{2}^{\downarrow}\right)$. To prove the proposition, it suffices to show that

$$
\begin{equation*}
f\left(\sigma_{1}^{\downarrow}, \sigma_{2}^{\downarrow}\right)=\max _{\rho_{1}^{\downarrow}} f\left(\rho_{1}^{\downarrow}, \sigma_{2}^{\downarrow}\right)=\min _{\rho_{2}^{\downarrow}} f\left(\sigma_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right) . \tag{B.29}
\end{equation*}
$$

Indeed, this is because B.29 enables the following calculation to go through:

$$
\begin{align*}
u_{1}(\sigma) & =f\left(\sigma_{1}^{\downarrow}, \sigma_{2}^{\downarrow}\right)=\max _{\rho_{1}^{\downarrow}} f\left(\rho_{1}^{\downarrow}, \sigma_{2}^{\downarrow}\right)  \tag{B.30}\\
& \geq \min _{\rho_{2}^{\downarrow}} \max _{\rho_{1}^{\downarrow}} f\left(\rho_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right)=\operatorname{gv}_{p}\left(G\left(\sigma^{\mathcal{T}}\right)\right)=\max _{\rho_{1}^{\downarrow}} \min _{\rho_{2}^{\downarrow}} f\left(\rho_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right)  \tag{B.31}\\
& \geq \min _{\rho_{2}^{\downarrow}} f\left(\sigma_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right)=f\left(\sigma_{1}^{\downarrow}, \sigma_{2}^{\downarrow}\right)=u_{1}(\sigma) . \tag{B.32}
\end{align*}
$$

Finally, we prove the max part of B.29). (The proof of the min part is analogous.) By Theorem 1 and reachable optimality, we have

$$
\begin{equation*}
u_{p}(\sigma)=\sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\sigma}(I) V_{p}^{\sigma}(I)=\sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\sigma}(I) \max _{\rho_{p} \supset \sigma_{p}^{\mathcal{T}}} V_{p}^{\rho_{p}, \sigma_{-p}}(I) \tag{B.33}
\end{equation*}
$$

Since the infoset values $V^{\tilde{\sigma}_{p}, \sigma_{-p}}(I)$ do not depend on the trunk portion of $p$ 's strategy and can be separately optimized in each $I$, the maximum and summation can be swapped, which concludes the proof:

$$
\begin{equation*}
f\left(\sigma^{\downarrow}\right)=u_{1}(\sigma)=\max _{\rho_{1} \supset \sigma_{1}^{\mathcal{T}}} \sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\sigma}(I) V^{\rho_{1}, \sigma_{2}}(I)=\max _{\rho_{2}^{\downarrow}} f\left(\sigma_{1}^{\downarrow}, \rho_{2}^{\downarrow}\right) \tag{B.34}
\end{equation*}
$$

Proposition 4.19 (Computing counterfactually-optimal values). Suppose that for each $\sigma^{\mathcal{T}}$, a value function $\mathbf{v}$ is of the form $\mathbf{v}^{\sigma^{\mathcal{T}}}(h)=v_{1}^{\mu}(h)$, where $\mu \in \Sigma$ is obtained by

- starting with some $\sigma \in \operatorname{NE}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$,
- for both $p=1,2$, going through all $I \subset \mathcal{Z}^{\mathcal{T}}, I \in \mathcal{I}_{p}$, that are counterfactually reachable by $p$ but not reachable,
- and replacing $\sigma_{p}$ by cbr $r_{p}\left(\sigma_{-p}\right)$ on such infosets and their descendants.

Then $\mathbf{v}$ is counterfactually optimal.
Proof. Firstly, $V_{p}^{\mu_{p}, \sigma_{-p}}(I)=\max _{\mu_{p}^{\prime}} V_{p}^{\mu_{p}^{\prime}, \sigma_{-p}}(I)=: V_{p}^{*, \sigma_{-p}}(I)$ holds for all reachable $I \subset \mathcal{Z}^{\mathcal{T}}$, since $\mu_{p}$ and $\sigma_{p}$ coincide there, and $\sigma$ couldn't be a Nash equilibrium otherwise (as we have shown in detail in the proof of Proposition 4.16). Secondly, an elementary backwards-induction argument implies that all infosets $I \subset \mathcal{Z}^{\mathcal{T}}$ that are unreachable but counterfactually reachable by $p$ satisfy $V_{p}^{\mu_{p}, \sigma_{-p}}(I)=V_{p}^{*, \sigma_{-p}}(I)$. Taken together, the two observations imply that $V_{p}^{\mu_{p}, \sigma_{-p}}(I)=V_{p}^{*, \sigma_{-p}}(I)$ holds for all counterfactually reachable sets. To prove the proposition, we need to show that the same holds if $\sigma_{-p}$ is replaced by $\mu_{-p}$.

Let $I \subset \mathcal{Z}^{\mathcal{T}}, I \in \mathcal{I}_{p}$ be counterfactually reachable by $p$. Applying Lemma 3.3 to $V_{p}^{\rho_{p}, \sigma_{-p}}(I)$, we get $V_{p}^{\rho_{p}, \sigma_{-p}}(I)=\sum_{h \in I} \frac{P_{-p}^{\sigma \mathcal{T}}(h)}{P_{-p}^{\sigma \mathcal{T}}(I)} v_{p}^{\rho_{p}, \sigma_{-p}}(h)$ for any $\rho_{p} \in \Sigma_{p}$. The histories $h \in I$ can be divided into two parts: those for which $P_{-p}^{\sigma^{\mathcal{T}}}(h)$ is positive, and those for which it is zero. The above formula shows that $V_{-p}^{\rho_{p}, \sigma_{-p}}(I)$ depends on the values of histories of the first type. However, such histories will be either fully reachable (because $P_{p}^{\sigma^{\tau}}(h)>0$ ) or counterfactually-unreachable by $p$ (because $P_{p}^{\sigma^{\top}}(h)=0$ ). In either case, they do not fit the criterion "cf. reachable but not reachable" for the opponent from the assumptions of this proposition, so $-p$ wasn't allowed to change their strategy below them and we have $v_{p}^{\rho_{p}, \mu_{-p}}(h)=v_{p}^{\rho_{p}, \sigma_{-p}}(h)$. It follows that $V_{p}^{\rho_{p}, \sigma_{-p}}(I)=V_{p}^{\rho_{p}, \mu_{-p}}(I)$ for every strategy of $p$, which concludes the proof.

Proposition 4.20 (Enabling DL-CFR). Let ( $\mathcal{T}$, v) be a depth-limited game corresponding to a counterfactually-optimal value function. Then:
(1) $D L-C F R$ can be viewed an instance of $C F R-D$ and inherits its guarantees.
(2) In particular, the strategy $\bar{\sigma}^{\mathcal{T}, t}$ produced after titerations of $D L-C F R$ satisfies $\operatorname{expl}\left(\bar{\sigma}^{\mathcal{T}}, t\right) \xrightarrow{t \rightarrow \infty} 0$.
In general, the same is not true for reachably-optimal value functions.
Proof. The "reachably-optimal value functions" part holds by Example 4.18
Comparing DL-CFR to CFR-D ( 20 [Alg. 6]), we see that CFR-D is essentially a CFR run in the trunk (like DL-CFR) that additionally requires a "Solve Subgames" method. This method takes the trunk strategy as an input and returns counterfactual values of infosets in $\mathcal{Z}^{\mathcal{T}}$ as output. In DL-CFR, these can be obtained by calling the value function, summing the values over $I$, and weighting them by the counterfactual reach probability of $I$. This proves (1). Theorem 17 from [20] then guarantees that the resulting strategy is an Nash equilibrium if the counterfactual values $V_{\mathrm{cf}}^{(\cdot)}(I)=P_{-p}^{(\cdot)}(I) V^{(\cdot)}(I)$ are maximal for all $I \subset \mathcal{Z}^{\mathcal{T}}$ (in the sense of corresponding to a strategy from which $p$ doesn't want to deviate).

This is equivalent to $V_{p}^{(\cdot)}(I)$ being maximal in counterfactually-reachable infosets. In DL-CFR with a counterfactually-optimal value function, this holds trivially, since $\mathbf{v}$ corresponds to strategies that satisfy this by definition.

Proposition 4.22 (NE as mutual naive best-responses). If $\mathbf{v}$ is a a counterfactually optimal value function, then any $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ where each $\sigma_{p}^{\mathcal{T}}$ is a naive best-response to $\sigma^{\mathcal{T}}$ in $(\mathcal{T}, \mathbf{v})$ can be extended into a Nash equilibrium in $G$.
Proof. Let $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ be s.t. each $\sigma_{p}^{\mathcal{T}}$ is a naive best response to $\sigma^{\mathcal{T}}$ in ( $\left.\mathcal{T}, \mathbf{v}\right)$. Denote by $\sigma$ the counterfactually-optimal extension of $\sigma^{\mathcal{T}}$ that $\mathbf{v}$ corresponds to. We will show that $\sigma$ is a Nash equilibrium in $G$.

Suppose that $p$ is considering to switch over to a different strategy $\rho_{p} \in \Sigma_{p}$ in $G$. We then have

$$
\begin{align*}
u_{p}\left(\rho_{p}, \sigma_{-p}\right) & =\sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\rho_{p}, \sigma_{-p}}(I) V_{p}^{\rho_{p}, \sigma_{-p}}(I) \\
& \leq \sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\rho_{p}, \sigma_{-p}}(I) \max _{\rho_{p}^{\prime} \in \Sigma_{p}} V_{p}^{\rho_{p}^{\prime}, \sigma_{-p}}(I) \tag{B.35}
\end{align*}
$$

Since $\mathbf{v}$ is counterfactually optimal, the maximum in B.35 is attained by $\sigma_{p}$ whenever $P^{\rho_{p}, \sigma_{-p}}(I)$ is non-zero, allowing us to continue as follows:

$$
\begin{aligned}
& \leq \sum_{I \subset \mathcal{Z}^{\mathcal{T}}} P^{\rho_{p}, \sigma_{-p}}(I) V_{p}^{\sigma}(I)=\sum_{h \in \mathcal{Z}^{\mathcal{T}}} P^{\rho_{p}, \sigma_{-p}}(h) v_{p}^{\sigma}(h) \\
& =\sum_{h \in \mathcal{Z}^{\mathcal{T}}} P^{\rho_{p} \mid \mathcal{T}, \sigma_{-p}^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h) \leq \max _{\mu_{p}^{\mathcal{T}} \in \Sigma_{p}^{\mathcal{T}}} \sum_{h \in \mathcal{Z}^{\mathcal{T}}} P^{\mu_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}}(h) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(h),
\end{aligned}
$$

where the first identity also used the fact that either $P^{\rho_{p}, \sigma_{-p}}(I)=0$ or

$$
P^{\rho_{p}, \sigma_{-p}}(I) V_{p}^{\sigma}(I)=P_{p}^{\rho_{p}}(I) P_{-p}^{\sigma_{-p}}(I) \sum_{h \in I} \frac{P_{-p}^{\sigma_{-p}}(h)}{P_{-p}^{\sigma-p}(I)} v_{p}^{\sigma}(h)=\sum_{h \in I} P^{\rho_{p}, \sigma_{-p}}(h) v_{p}^{\sigma}(h) .
$$

Since $\sigma_{p}^{\mathcal{T}}$ is a naive best response to $\sigma^{\mathcal{T}}$, the last term is equal to $u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)$, which is further equal to $u_{p}(\sigma)$. We have shown that $p$ cannot increase their utility by deviating from $\sigma_{p}$, which concludes the proof.

Proposition 4.24 (Inspired by [4]). Suppose that $\mathbf{v}$ is reachabl $\sqrt{41}$ optimal w.r.t. a portfolio that contains (the trunk restrictions of) all pure undominated strategies. Then $\mathbf{v}$ preserves the equilibri ${ }^{42}$ of $G$.

[^22]Proof. By Theorem 3, the conclusion holds for reachably optimal value functions (those from Definition 4.15 i.e., those optimal w.r.t. all strategies). We will show that if $\mathbf{v}$ and $\mathbb{P}$ satisfy the assumptions of the theorem, $\mathbf{v}$ must be reachably optimal. To do this, it suffices to show that any $\sigma \supset \sigma^{\mathcal{T}}$ extension that is reachably optimal w.r.t. $\mathbb{P}$ is reachably optimal.

Suppose that $\sigma \supset \sigma^{\mathcal{T}}$ is not reachably optimal. By definition, this means that there is some trunk-leaf infoset $I$ with $P^{\sigma^{\mathcal{T}}}(I)>0$ and $p$ such that $V_{p}^{\rho_{p}, \sigma_{-p}}(I)>$ $V_{p}^{\sigma_{p}, \sigma_{-p}}(I)$. Without loss of generality, assume that $\rho_{p}$ is an undominated pure strategy. (If such $\rho_{p}$ exists within the space of all strategies, there will also exist some pure $\rho_{p}^{\prime}$ with the same property, and some undominated pure $\rho_{p}^{\prime \prime}$ with the same property.) Since $\left.\rho_{p}\right|_{\mathcal{H} \backslash \mathcal{T}} \in \mathbb{P}_{p}$, this shows that the extension $\sigma$ is not reachably optimal w.r.t. $\mathbb{P}$.

Theorem 4. For any $\mathcal{T}, \mathbb{P}$, and $\mathbf{v}$ that is reachably optimal w.r.t. $\mathbb{P}$, a trunk strategy $\sigma^{\mathcal{T}}$ is solution of $(\mathcal{T}, \mathbf{v})$ if and only it is in $\left.\operatorname{NE}(G(\mathcal{T}, \mathbb{P}))\right|_{\mathcal{T}}$.

Proof. First, we prove a lemma about the equality of the depth-limited expected utility and the game value of the game $G(\mathcal{T}, \mathbb{P})\left(\sigma^{\mathcal{T}}\right)$ where the trunk strategy is fixed to $\sigma^{\mathcal{T}}$ (so the players only select which portfolio strategy to use below each infoset) ${ }^{43}$
Lemma. If $\mathbf{v}$ is reachably optimal w.r.t. $\mathbb{P}$, we have $u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\operatorname{gv}_{p}\left(G(\mathcal{T}, \mathbb{P})\left(\sigma^{\mathcal{T}}\right)\right)$ for every $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$.

To prove the lemma, suppose that $\mathbf{v}$ is reachably optimal w.r.t. $\mathbb{P}$ and let $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$. By Definition 4.23, there is some extension $\sigma$ of $\sigma^{\mathcal{T}}$ which is reachably optimal w.r.t. $\mathbb{P}$ and satisfies $v_{p}^{\sigma}(z)=\mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(z)$ for all reachable $z \in \mathcal{Z}^{\mathcal{T}}$. Since only the reachable $z$ make a difference when calculating utility, we have $u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)=\sum_{z \in \mathcal{Z}^{\mathcal{T}}} P^{\sigma^{\mathcal{T}}}(z) \mathbf{v}_{p}^{\sigma^{\mathcal{T}}}(z)=\sum_{z \in \mathcal{Z} \mathcal{T}} P^{\sigma}(z) v_{p}^{\sigma}(z)=u_{p}(\sigma)$. By $\sigma^{\mathbb{P}}$, we denote the strategy from the game $G(\mathcal{T}, \mathbb{P})$ which corresponds to $\sigma$ (formally, $\sigma^{\mathbb{P}}$ coincides with $\sigma^{\mathcal{T}}$ on $\mathcal{T}$, it is defined by condition (i) from Definition 4.23 in reachable $I \subset \mathcal{Z}^{\mathcal{T}}$, and it is defined arbitrarily for unreachable $I$ ). From definition of $G(\mathcal{T}, \mathbb{P})$, it follows that $u_{p}(\sigma)$ is equal the expected utility of $\sigma^{\mathbb{P}}$ in $G(\mathcal{T}, \mathbb{P})$.

To finish the proof of the lemma, it remains to show that $\sigma^{\mathbb{P}}$ is an NE of $G(\mathcal{T}, \mathbb{P})\left(\sigma^{\mathcal{T}}\right)$ (and hence its expected utility is equal to the value of this game). To see how the converse would lead to a contradiction, note that if $\sigma^{\mathbb{P}}$ wasn't an equilibrium, one player $q$ would be able to increase their utility in $G(\mathcal{T}, \mathbb{P})$ by deviating from $\sigma^{\mathbb{P}}$. Since $\sigma^{\mathcal{T}}$ is fixed, this deviation would need to happen in some $I \subset \mathcal{Z}^{\mathcal{T}}$. Moreover, to affect the total expected utility, $I$ would need to be reachable. And without loss of generality, the new strategy in $G(\mathcal{T}, \mathbb{P})$ would be pure - in $G$, this would translate to a deterministic choice of some $\rho_{q}^{\downarrow} \in \mathbb{P}_{q}$. We

[^23]would then have $V^{\rho_{q}^{\downarrow}, \sigma_{-p}}(I)>V^{\sigma}(I)$ - a contradiction with the assumption that $\sigma$ is reachably optimal w.r.t. $\mathbb{P}$.

We now finish the proof of the theorem. Let $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$. We will show that $\sigma^{\mathcal{T}}$ is not a solution of $(\mathcal{T}, \mathbf{v})$ if and only if it is not a restriction of some NE of $G(\mathcal{T}, \mathbb{P})$. Using the lemma, we see that (A) the existence of some $\rho_{p}^{\mathcal{T}}$ for which $u_{p}^{\mathbf{v}}\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)>u_{p}^{\mathbf{v}}\left(\sigma^{\mathcal{T}}\right)$ is equivalent to (B) the existence of some (or rather, the same one) $\rho_{p}^{\mathcal{T}}$ for which $\operatorname{gv}_{p}\left(G(\mathcal{T}, \mathbb{P})\left(\rho_{p}^{\mathcal{T}}, \sigma_{-p}^{\mathcal{T}}\right)\right)>\operatorname{gv}_{p}\left(G(\mathcal{T}, \mathbb{P})\left(\sigma^{\mathcal{T}}\right)\right)$. By Proposition 4.5. (B) is equivalent to $\sigma^{\mathcal{T}}$ being exploitable and thus (by Lemma 4.4) not being a restriction of some Nash equilibrium of $G(\mathcal{T}, \mathbb{P})$. Since (A) is literally the definition of $\sigma^{\mathcal{T}}$ not being the solution $(\mathcal{T}, \mathbf{v})$, this concludes the proof.

Proposition 4.27 (Sufficient statistics for optimal value functions). Let $\mathcal{T}$ be a trunk and $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$.
(1) The (joint) reach probabilities $\left(P^{\sigma^{\tau}}(h)\right)_{h \in \mathcal{Z}^{\mathcal{T}}}$ provide a sufficient statistic for computing some reachably-optimal $\mathbf{v}$.
(2) The factored reach probabilities $\left(P_{1}^{\sigma^{\top}}(h), P_{2}^{\sigma^{\top}}(h)\right)_{h \in \mathcal{Z} \mathcal{T}}$ provide a sufficient statistic for computing some counterfactually-optimal $\mathbf{v}$.
In particular, it suffices to keep $\left(P_{1}^{\sigma^{\mathcal{T}}}(I)\right)_{\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{1}}$ and $\left(P_{2}^{\sigma^{\mathcal{T}}}(I)\right)_{\mathcal{Z}^{\mathcal{T}} \supset I \in \mathcal{I}_{2}}$.
Proof. The first case follows from Proposition 4.16. Indeed, the proposition states that a reachably-optimal value function can be found by identifying a solution of $\sigma \in G\left(\sigma^{\mathcal{T}}\right)$ and computing the corresponding values $v_{1}^{\sigma}(h), h \in \mathcal{Z}^{\mathcal{T}}$. To calculate $v_{1}^{\sigma}(h)$, we only need to know how $\sigma$ looks in the bottom of the game (i.e., its trunk-portion isn't needed). To find this bottom part of $\sigma$, recall that $G\left(\sigma^{\mathcal{T}}\right)$ looks like $G$, except that all decisions in $\mathcal{T}$ are done by chance, according to the probabilities that $\sigma^{\mathcal{T}}$ prescribes for the given state. We can replace the whole trunk by a single chance node, where the probability of transitioning to $h \in \mathcal{Z}^{\mathcal{T}}$ is $P^{\sigma^{\mathcal{T}}}(h)$. This modified game only requires the knowledge of the reach probabilities over $\mathcal{Z}^{\mathcal{T}}$ while still being able to recover $\left.\sigma\right|_{\mathcal{H} \backslash \mathcal{T}}$.

The second case follows from Proposition 4.19. Indeed, the "in particular" part is true because $P_{p}^{\sigma^{\mathcal{T}}}(I)=P_{p}^{\sigma^{\mathcal{T}}}(h)$ holds whenever $h \in I \in \mathcal{I}_{p}$. In turn, the separated reach probabilities of individual histories allow us to recover the joint probabilities via $P^{\sigma^{\mathcal{T}}}(h)=P_{1}^{\sigma^{\mathcal{T}}}(h) P_{2}^{\sigma^{\mathcal{T}}}(h) P_{c}(h)$. By the first part of this proposition, this enables us to find some (bottom part of) $\sigma \in \mathrm{NE}\left(G\left(\sigma^{\mathcal{T}}\right)\right)$. To perform the post-processing step, we need to know the beliefs $P^{\sigma}(\cdot \mid J)$ for all counterfactually reachable infosets $I \subset \mathcal{H} \backslash \mathcal{T}$. However, these can all be calculated by starting with the beliefs $P^{\sigma^{\mathcal{T}}}(\cdot \mid I), I \subset \mathcal{Z}^{\mathcal{T}}$, and using $\left.\sigma\right|_{\mathcal{H} \backslash \mathcal{T}}$. The beliefs $P^{\sigma^{\tau}}(\cdot \mid I)$ are only known for counterfactually reachable infosets however, they are also only needed for such sets (since the post-processing is only needed in counterfactually reachable $J$ ).

Proposition 4.28 (Localization by public states). For any public state $S \subset \mathcal{Z}^{\mathcal{T}}$ :
(i) Both $\left(P^{\sigma^{\top}}(h)\right)_{h \in S}$ and $\left(P^{\sigma^{\top}}(h \mid S)\right)_{h \in S}$ provide a sufficient statistic for computing some reachably-optimal value function $\mathbf{v}^{\sigma^{\top}}(\cdot)$ on $S$.
(ii) $\left(P_{1}^{\sigma^{\mathcal{T}}}(I)\right)_{S \supset I \in \mathcal{I}_{1}}$ and $\left(P_{2}^{\sigma^{\mathcal{T}}}(I)\right)_{S \supset I \in \mathcal{I}_{2}}$ together provide a sufficient statistic for computing some counterfactually-optimal $\mathbf{v}^{\sigma^{\tau}}(\cdot)$ on $S$.

Proof. Suppose we are given some statistic $X: \Sigma^{\mathcal{T}} \times\left\{S \subset \mathcal{Z}^{\mathcal{T}} \mid S \in \mathcal{S}\right\} \rightarrow \mathcal{X}$ as described by the assumptions of one of the cases of the proposition (e.g., $\left.X\left(\sigma^{\mathcal{T}}, S\right)=\left(P^{\sigma^{\mathcal{T}}}(h)\right)_{h \in S}\right)$. To prove the result, we need to construct a value function that is optimal in the appropriate sense and show that there exists a function $\tilde{\mathbf{v}}: \mathcal{Z}^{\mathcal{T}} \times \mathcal{X} \rightarrow \mathbb{R}$ for which $\tilde{\mathbf{v}}^{X\left(\sigma^{\mathcal{T}}, S(h)\right)}(h)=\mathbf{v}^{\sigma^{\mathcal{T}}}(h)$ (where $S(h)$ denotes the public state that $h$ belongs to).

Without loss of generality, suppose that the function $X$ is surjective. For $S \subset \mathcal{Z}^{\mathcal{T}}$ and $X \in \mathcal{X}$, pick an arbitrary $\sigma^{\mathcal{T}} \in \Sigma^{\mathcal{T}}$ with $X\left(\sigma^{\mathcal{T}}, S\right)=X$ and denote by $\sigma$ its extenstion that is optimal in the appropriate sense. (By Proposition 4.16 resp. 4.19, this is always possible.) Denote by $\rho^{S, X}$ the restriction of $\sigma$ to $G(S)$.

Observe that while there was a lot of ambiguity in the choice of $\sigma^{\mathcal{T}}$ and $\sigma$, any extension $\mu \in \Sigma$ of $\rho^{S, X}$ that satisfies $X\left(\left.\mu\right|_{\mathcal{T}}, S\right)=X$ will also satisfy the appropriate definition of optimality for $\left.\mu\right|_{\mathcal{T}}$ on $S$. Indeed, this is true because
(i) $\sigma$ satisfies the appropriate definition for $\left.\sigma\right|_{\mathcal{T}}$ and $S$ and
(ii) all the terms present in the corresponding definition only depend on $\left.\mu\right|_{G(S)}$ and $X\left(\left.\mu\right|_{\mathcal{T}}, S\right)$, which coincide with $\left.\sigma\right|_{G(S)}$ and $X\left(\left.\sigma\right|_{\mathcal{T}}, S\right)$.
(This is the crucial part of the proof which only goes through because $S$ is closed under infosets - without this, we might not have all information necessary to compute $V_{p}^{(\cdot)}(I)$.)

For $\sigma^{\mathcal{T}}$, define $\mu\left(\sigma^{\mathcal{T}}\right):=\sigma^{\mathcal{T}} \cup \bigcup_{S \subset \mathcal{Z} \mathcal{T}} \rho^{S, X\left(\sigma^{\mathcal{T}}, S\right)}$ and $\mathbf{v}^{\sigma^{\mathcal{T}}}(h):=\mathbf{v}_{1}^{\mu\left(\sigma^{\mathcal{T}}\right)}(h)$. By the above observation, $\mathbf{v}$ is a value function that is optimal in the appropriate sense.

Finally, denote by $\rho(S, X)$ an arbitrarily chosen extension of $\rho^{S, X}$ which satisfies $X\left(\left.\rho(S, X)\right|_{\mathcal{T}}, S\right)=X$ and define $\tilde{\mathbf{v}}^{X}(h):=\mathbf{v}_{1}^{\rho(S(h), X)}(h)$. To verify the definition of $X$ being a sufficient statistic for $\mathbf{v}$, let $\sigma^{\mathcal{T}}$ be a trunk strategy $\sigma^{\mathcal{T}}$ and $h \in S \subset \mathcal{Z}^{\mathcal{T}}$. By their definitions, the strategies $\rho\left(S, X\left(\sigma^{\mathcal{T}}, S\right)\right)$ and $\mu\left(\sigma^{\mathcal{T}}\right)$ coincide on $G(S)$. In particular, they coincide on all descendants of $h$ and we have

$$
\mathbf{v}^{\sigma^{\mathcal{T}}}(h)=\mathbf{v}_{1}^{\mu\left(\sigma^{\mathcal{T}}\right)}(h)=\mathbf{v}_{1}^{\rho\left(S, X\left(\sigma^{\mathcal{T}}, S\right)\right)}(h)=\tilde{\mathbf{v}}^{X\left(\sigma^{\mathcal{T}}, S\right)}(h)
$$

which concludes the proof.
Theorem 5 (Public state minimality). Let $\mathcal{T}$ be a trunk, $\mathcal{Z}^{\mathcal{T}}$ its leaves, $S \subset$ $\mathcal{Z}^{\mathcal{T}}$ a common-knowledge public state, and $h_{0}, g \in S$. Suppose that trunk strategies $\sigma^{\mathcal{T}}$ and $\mu^{\mathcal{T}}$ render the same non-zero reach probabilities at $S$, except that $P^{\sigma^{\mathcal{T}}}(g) \neq P^{\mu^{\top}}(g)=0$.

Then there exists some game $G$, s.t. $\mathcal{T}$ is a trunk in $G$ and $\mathcal{Z}^{\mathcal{T}}$ the corresponding leaves, for which both $G\left(\sigma^{\mathcal{T}}\right)$ and $G\left(\mu^{\mathcal{T}}\right)$ each have a unique Nash equilibrium, $\sigma$ and $\mu$, and these satisfy $v_{1}^{\sigma}\left(h_{0}\right) \neq v_{1}^{\mu}\left(h_{0}\right)$.

Proof. Let $\mathcal{T}, \mathcal{Z}^{\mathcal{T}}, S, h_{0}, g, \sigma^{\mathcal{T}}$ and $\mu^{\mathcal{T}}$ be as in Theorem 5 . We shall prove the theorem by constructing $G$ and showing that it has the desired properties.

We start by making two simplifying assumptions. First, we assume that each $h \in S$ only has one legal (dummy) action that we denote $d$. In the general case, each $h a, h b$ would be extended identically, complicating the notation but not introducing any real challenges. Since the public state cannot be further refined, there exists a sequence of histories satisfying $h_{0} \sim h_{1} \sim \cdots \sim h_{n}=g$ in $S$, where $\sim$ denotes "one of the players cannot distinguish between the two histories". We assume, without loss of generality, that $h_{0}, \ldots, h_{n}$ is the shortest among such sequences. In particular, it follows that the sequence looks like the one in Figure 2 (histories and information sets don't repeat, players unable to distinguish them alternate). We only show the proof in the case where both $h_{0}, h_{1}$ and $h_{n-1}, h_{n}$ are indistinguishable by the first player (the proofs of the remaining three cases are similar). (Note that the proof only requires the reach probabilities to be non-zero for the histories $h_{i}$, rather than on the whole $S$.)

The history tree $\mathcal{H}$ of $G$ is s.t. $\mathcal{T}$ is a trunk in $G$ and $\mathcal{Z}^{\mathcal{T}}$ its leaves. Outside of $S, G$ can continue arbitrarily - e.g., by each node $h \in \mathcal{Z}^{\mathcal{T}} \backslash S$ being terminal with $u_{1}(h)=0$. Finally, at $\{h d \mid h \in S\}$ and below, $G$ is defined as follows :

- For $h \in S \backslash\left\{h_{i} \mid i=0, \ldots, n\right\}, h d$ is a terminal node with utility 0 (for all strategic considerations, this replaces $S$ by $\left.\left\{h_{i} \mid i=0, \ldots, n\right\}\right)$.
- For $h_{i}, i=0, \ldots, n-1, h_{i} d$ leads to a matching pennies game (a matrix game with actions $U, D$ for player 1 , actions $L, R$ for player 2 , and corresponding utilities 1 for $U, L$ and $D, R$, resp. 0 for $U, R$ and $D, L)$.
- For $h_{n}, h_{n} d$ leads to a game where only player 1 acts, choosing between $U$ (utility 0 ) and $D$ (utility 1 ).
- The information sets below $h_{i}$ are defined in such a way that player 1 has to use the same strategy below $h_{0}$ and $h_{1}, h_{2}$ and $h_{3}, \ldots, h_{n-1}$ and $h_{n}$, and player 2 has to use the same strategy below $h_{1}$ and $h_{2}, h_{3}$ and $h_{4}, \ldots$, $h_{n-2}$ and $h_{n-1}$ (player 2 strategy below $h_{0}$ is independent of everything else).

Since $h_{n}$ is unreachable under $\mu^{\mathcal{T}}, \widetilde{G}\left(\mu^{\mathcal{T}}\right)$ below $S$ is effectively a collection of (interconnected) matching pennies games. It follows that the uniform strategy of both players is a Nash equilibrium (clearly, no player can improve his overall utility). On the other hand, it is not a NE strategy in $\widetilde{G}\left(\sigma^{\mathcal{T}}\right)$ below $S$ (since player 1 could improve his utility be deviating to " $D$ everywhere").

In particular, 0.5 is an expected utility of $h_{0}$ under some NE strategy in $G\left(\mu^{\mathcal{T}}\right)$. Suppose that some NE strategy $\sigma$ in $G\left(\sigma^{\mathcal{T}}\right)$ has $v_{1}^{\sigma}\left(h_{0}\right)=0.5$. We will show that such $\sigma$ has to be uniformly random, and thus prove the theorem by contradiction.

Firstly, if $\sigma_{1}$ wasn't uniformly random at the information set $\left\{h_{0} d, h_{1} d\right\}$, player 2 could increase his overall utility by changing his strategy below $h_{0}$ to either $L$ or $R$ and thus $\sigma$ wouldn't be a NE.

We proceed inductively. We know that $u_{1}^{\sigma}\left(h_{0}\right)=0.5$, and that in $\left\{h_{0} d, h_{1} d\right\}$, $\sigma_{1}$ takes both $U$ and $D$ with non-zero probability. If $\sigma_{1}$ is to be a NE, player 1 has to be indifferent between playing $U$ and $D$ in $\left\{h_{0} d, h_{1} d\right\}$. Since $v_{1}^{\sigma}\left(h_{0}\right)=0.5$, this can only be achieved if $\sigma_{2}$ takes both $L$ and $R$ below $h_{1}$ with the same probability. In particular, $u_{1}^{\sigma}\left(h_{1}\right)=0.5$. Since $G$ forces the strategy of player 2 to be the same below $h_{1}$ and $h_{2}$, we get that $\sigma_{2}$ is uniformly random below $h_{2}$ as well.

We repeat the argument above for each $h_{i}$, eventually showing that if the players are to be indifferent between the actions they take with non-zero probability, $\sigma_{1}$ has to be uniformly random in the whole $G(S)$ and $\sigma_{2}$ has to be uniformly random below $h_{1}, \ldots, h_{n-1}$. Finally, if $\sigma_{1}$ wasn't uniform below $h_{0}$ (but was below $h_{1}$ ), player 1 could increase his utility by deviating to either $U$ or $D$. This implies that the whole $\sigma$ is uniformly random, which contradicts our earlier observation.

## C. Detailed Descriptions of the Domains

In this section, we give a more specific description of the domains used for evaluation.

## C.1. The Rules of Leduc hold'em ( $\mathbf{L H \text { ) }}$

Leduc Hold'em, which was first introduced in [38], is played with six cards and two suits: Two Jacks, two Queens, and two Kings; on of each per suit. Each player gets dealt a card and submits an ante of 1 . Each betting round has a maximum of two raise actions. (In other words, the possible actions are either raise or check if acting first in a betting round, re-raise, call or fold if faced with a raise, and call or fold if faced with a re-raise.) There are two betting rounds. After the first betting round, a public card will be shown and a second second betting round is performed using the same rules. The raise size is 2 in the first betting round and 4 in the second. A player wins if their private card matches the public card and the opponent's card does not match, or if none of the players' cards matches the public card and their private card is higher than the private card of the opponent. If no player wins, the game is a draw, and the pot is split. The maximum utility in this game is 13 (but we normalize it to 1 in the experiments).

## C.2. The Rules of Imperfect-Information Goofspiel (GS)

In goofspiel with $N$ cards, each player is given a private hand of bid cards with values 1 to $N$. A different deck of $N$ point cards is placed face up in a stack. Each turn, one of the point cards is revealed and each player bids for it by secretly choosing a single card in their hand. The highest bidder gets the point card and adds the point total to their score (in case of a tie, the point card is discarded). This is repeated $N$-times until both players run out of cards, and the player with the highest score wins. In imperfect-information goofspiel, the players only discover who won or lost a bid, but not which bid cards were chosen.

We assume that the point cards are revealed in a decreasing order and that this is common knowledge. We use $N=5$ and define the each player's utilities as the sum of the cards they won minus the sum of the opponent's cards. As a result, the maximum utility is $(5+4+3+2)-1=13$ (obtained when the opponent only wins the last card).Note that the structure of the game is different from Leduc hold'em:First, the actions are not perfectly observable. Second, if a player plays a card and loses, the opponent might have played any higher card - this implies that the size of the information sets and public states first grows at the start of the game and then shrinks towards the end

## C.3. The Rules of Oshi-zumo (OZ)

Oshi-zumo [48] is played by two players, both of whom start with $N$ coins. At the beginning of the game, a sumo wrestler is positioned at the center of a one-dimensional playing field which consists of $2 K+1$ locations. Each round each player secretly spends between some number of coins - at least $M$ and at most the number of coins they have remaining. (If a player ends up with no coins, while their opponent still has some, the opponent immediately wins.) Whoever spent more coins then gets to push the wrestler one location towards the opponent's side. If this causes the wrestler to fall off the playing field, the player who did the pushing wins. (If the bids are equal, the wrestler does not move.) We use an imperfect-information variant of the game, where the winner is publicly revealed, but it is not revealed how many coins each of the players spent. If both players run out of coins at the same time, the player who pushed the wrestler further (i.e., the player who does not have the wrestler on their side) wins. (If the wrestler is located at the center, the game result is a draw.) In our variant of the game, winning yields 1 utility. We use $N=8$ and $K=1$.

Like goofspiel, oshi-zumo is comprised of variably-sized information sets resulting from the uncertainty about the opponent's bid. Unlike in GS, not all terminal branches in OZ have the same length (because of the possibility to run out of coins prematurely).

## D. Choice of the Loss Function

A crucial part of approximating value functions is to understand which loss to minimize to achieve the best performance.

Minimizing Huber, $l_{1}$ or $l_{\infty}$. To explore and determine which of the standard losses are best, we used both the same architecture and data to train value networks, validating on all three candidate loss functions, however each of them minimizing either Huber, $l_{1}$ or $l_{\infty}$ loss. Across the three cases, we then compared the minimum validation losses which could be achieved under the particular setting. This analysis revealed that there is no significant difference between all three losses in terms of their achieved validation losses.

We also performed an extensive analysis of alternative loss functions Below, we describe the high-level ideas.

Range-weighting. Since the main metric, we care about in the value network, is the achieved exploitability when used in conjunction with $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$, an intuitive idea would be to prioritize the error on infosets which are actually played to under a certain trunk strategy, i.e., the range of the infoset. If an arbitrary infoset's range is zero, the value network's error on it should be less important than one which has high range.Hence, we implemented a loss function which computes the standard Huber, $l_{1}$ or $l_{\infty}$ error on every sample, however multiplies the error by the particular range of the infoset.

Public state value weighting. A similar idea was to emphasize errors not just on particular infosets but on full public states, i.e., a full sample. As some public states have higher value than others, doing badly in them might hurt the resulting exploitability more than in low value public states. We did so by precomputing the values of each public state and then multiplying the resulting errors with these values.

Zero-sum Loss. Another idea, already realized by [1] was the zero-sum loss which is computed by multiplying each player's ranges with their corresponding counterfactual values and summing the products up to a scalar for each player. The sum of both of these numbers should then be equal to 0 . This loss can then be minimized in conjunction with standard losses.

Enforcing the right ratio between values. A general phenomenon we observed when analysing the game-solving behavior of $\mathrm{DL}-\mathrm{CFR}_{\mathrm{NN}}^{+}$was also the algorithm's vulnerability in regard to values which should be negative or exact zero but are, due to naturally occurring neural network inaccuracies, very small positive numbers. Relatedly, in regret matching absolute values of regrets do not really matter, but rather the ratio between all values. In the case where an infoset's correct regret matching strategy should be pure for one action this problem can change a strategy from pure to mixed and cause the algorithm to visit states which result in high exploitability. This leads us to conclude that predicting the right ratio between counterfactual values is key and should be enforced through the loss signal. Since negative values in regret matching are pruned to 0 , we implemented a loss function which first shifts all values into the positive range using a constant and then applies softmax to both the targets and predictions and minimizes the Kullback-Leibler divergence between the two.

Summary of explored loss functions. All of the tried variations did not improve (1) the speed of convergence (2) proved to be linear transformations of standard loss functions and, most importantly, did not have a stronger correlation with the resulting exploitability. Despite that, we decided to report the findings of this extensive search and concluded that using one of the standard losses is sufficient. Finally, we chose Huber as the loss to run all experiments on. The code for all above described loss functions is accessible under [49].

## E. Detailed Investigation of Value Networks

## E.1. Hyperparameter Optimization

We start out by describing the hyperparameter optimization performed in order to train our neural network. In particular, we try to find the right choice of network width and depth, loss function, and the amount of data. While our main results are presented using normalized losses, we opted to show non-normalized losses in this section. The reason for that is that the objective of this investigation is strictly about finding optimal settings for certain hyperparameters as opposed to judging certain performance criteria of the full algorithm across domains.

## E.1.1. Layer Width

The three following figures show the search through different widths of hidden layers for each domain. The y-axis denotes the best non-normalized error achieved with a given setting after 1000 epochs. The x -axis denotes the ratio of the number of neurons to the input size.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-74.jpg?height=619&width=1026&top_left_y=1160&top_left_x=517)

Figure E.7: Layer width grid search in Leduc hold'em. The y-axis denotes validation error while the x -axis shows number of neurons as a multiplier of the input size.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-75.jpg?height=614&width=1012&top_left_y=449&top_left_x=540)

Figure E.8: Layer width grid search in goofspiel. The y-axis denotes validation error while the x-axis shows number of neurons as a multiplier of the input size.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-75.jpg?height=625&width=1041&top_left_y=1303&top_left_x=504)

Figure E.9: Layer width grid search in oshi-zumo. The y-axis denotes validation error while the x -axis shows number of neurons as a multiplier of the input size.

## E.1.2. Network Depth

We also explored different neural network depths for each domain. We display the best non-normalized loss achieved (the y-axis) with the given number of layers (the x-axis) after 1000 epochs.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-76.jpg?height=602&width=1174&top_left_y=420&top_left_x=481)

Figure E.10: Layer depth grid search in Leduc hold'em. The y-axis denotes validation error while the x -axis shows number of hidden layers.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-76.jpg?height=600&width=1161&top_left_y=1207&top_left_x=498)

Figure E.11: Layer depth grid search in goofspiel. The y-axis denotes validation error while the x -axis shows number of hidden layers.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-77.jpg?height=608&width=1202&top_left_y=414&top_left_x=459)

Figure E.12: Layer depth grid search in oshi-zumo. The y-axis denotes validation error while the x -axis shows number of hidden layers.

Recall, that Table 2 in the main text shows the final architectures used, which were based on the above described grid searches.

## E.1.3. Data Amount

We are also interested in how much data is required to achieve a sufficiently low loss. Recall that the network always takes some public state $S$ and a range $r$ at $S$ as input (and returns a vector $\vec{v}$ of values, one for each infoset at $S$ ). We thus measure the amount of data in the number of $(S, r, \vec{v})$ tuples used for training. Note that each of the randomly generated training trunk-strategies (Section 6.1.2) corresponds to a as many tuples as there are public states at the depth limit (Table 3). Hence, the more public states there are in a domain the less random trunk strategies needed. We show the validation loss after training 1000 epochs for a given data amount setting (LH in Fig. E.13. GS in Fig. E. 14 , and OZ in Fig. E.15).
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-78.jpg?height=627&width=1018&top_left_y=445&top_left_x=532)

Figure E.13: Relationship of training samples on the validation errors in LH. We show loss on the $y$-axis and number of training samples on the x -axis.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-78.jpg?height=627&width=1015&top_left_y=1345&top_left_x=533)

Figure E.14: Relationship of training samples on the validation errors in GS. We show loss on the y -axis and number of training samples on the x -axis.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-79.jpg?height=616&width=1009&top_left_y=451&top_left_x=539)

Figure E.15: Relationship of training samples on the validation errors in OZ. We show loss on the y -axis and number of training samples on the x -axis.

We conclude that even with below 10k datapoints we can achieve low validation errors depending in every domain. However, it is clear that with more data it is possible to reach orders of magnitude better validation losses. This is however not necessary to guarantee low exploitability. Regardless of that, to run our main experiments, we opted to use a larger amount of data to properly be able to explore the algorithm's performance: the numbers we used for our main experiments were 18000 for goofspiel, 70644 for Leduc and 34000 for oshi-zumo.

## E.2. Performance on $C F R-D$ Data

We also wanted to see the performance of the network on the ranges which CFR-D requests when solving the game. To achieve that, we trained the network on the standard random ranges and validated on the CFR-D ranges. The following figures summarize this experiment for all domains.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-80.jpg?height=600&width=1006&top_left_y=426&top_left_x=538)

Figure E.16: LH: Training a network on random ranges and measuring validation error on CFR-D data.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-80.jpg?height=603&width=1004&top_left_y=1235&top_left_x=539)

Figure E.17: GS: Training a network on random ranges and measuring validation error on CFR-D data.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-81.jpg?height=597&width=1009&top_left_y=428&top_left_x=539)

Figure E.18: OZ: Training a network on random ranges and measuring validation error on CFR-D data.

We observe substantially higher $l_{\infty}$ and $l_{1}$-errors on CFR-D data compared to random ranges, while the Huber errors are close to the one on random ranges. We conclude that the CFR-D data contains problematic ranges which are not covered by the training data. Despite that, it is clear (1) that the network can generalize to CFR-D data sufficiently precisely on average and (2) that the validation loss on unseen random ranges and provides enough confidence about the resulting exploitability as shown in Figure 3 .

## E.3. Comparison of $D L-C F R_{N N}^{+}$and $C F R-D$ \& Training Ranges

Here, we show the remaining results for the two other domains. Recall, that we examined the game-solving behaviour of $\mathrm{DL}^{-} \mathrm{CFR}_{\mathrm{NN}}^{+}$in a particular public state (meaning the which infosets are actually played into, i.e. ranges) and compared it to CFR-D in the same iteration and the closest (in terms of Euclidean distance) random range in the training data set for that particular range of $\mathrm{DL}_{\mathrm{CFR}}^{\mathrm{NN}}+\mathrm{C}$. We show the ranges in the top row and the corresponding values, i.e., solutions to the subgames, in the bottom.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-82.jpg?height=375&width=1191&top_left_y=428&top_left_x=467)

Figure E.19: We show ranges (top row) and values (bottom row) of the first 20 iterations of CFR-D (left), DL-CFR ${ }_{\mathrm{NN}}^{+}$(middle) and closest training ranges (right) in goofspiel. The x-axis refers to a information set in a particular public state. The $y$-axis denotes the range/value in a given iteration.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-82.jpg?height=367&width=1186&top_left_y=1063&top_left_x=464)

Figure E.20: We show ranges (top row) and values (bottom row) of the first 10 iterations of CFR-D (left), $\mathrm{DL}^{-C F R}{ }_{\mathrm{NN}}^{+}$(middle) and closest training ranges (right) in oshi-zumo. The x-axis refers to a information set in a particular public state. The y-axis denotes the range/value in a given iteration.

We observe very similar ranges requested by CFR-D and DL-CFR ${ }_{\mathrm{NN}}^{+}$. The closest training ranges in all domains are still substantially different from DL-CFR ${ }_{\mathrm{NN}}^{+}$. However, in goofspiel, the corresponding predicted values of DL-CFR ${ }_{\mathrm{NN}}^{+}$are different from those computed by CFR-D, hinting at the possibility that $\mathrm{DL}_{\mathrm{CFR}}^{\mathrm{NN}}+\mathrm{was}$ able to identify a different trunk equilibrium. Note that the value networks which were used to generate this data all achieved near-optimal exploitability.

## E.4. Public State Cross-Validation

This subsection presents the results of the other two domains of the public state cross-validation experiment. In all three figures below the x -axis denotes the public state index, while the $y$-axis denotes the $l_{\infty}$ validation loss for a network which has not seen the particular public state in training, the validation loss of a fully trained network and the average of 10 random networks.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-83.jpg?height=418&width=787&top_left_y=428&top_left_x=452)

Figure E.21: Public state cross-validation in goofspiel.
![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-83.jpg?height=440&width=847&top_left_y=994&top_left_x=452)

Figure E.22: Public state cross-validation in oshi-zumo.
We see that the network with withheld data is nearly as bad as the random network in goofspiel, somewhere in the middle in oshi-zumo, and nearly as good as the full value network in LH. This is consistent with our hypothesis that the encoding is suitable for generalization: Goofspiel has a very small number of public states, so there is not much to generalize from. On the other hand, LH has many public states and there is a good chance that for any given $S$, there will be other public states that are strategically similar to $S$ - and indeed, the network generalizes very well here.


[^0]:    *To whom correspondence should be addressed.
    Email address: viliam.lisy@agents.fel.cvut.cz (Viliam Lisý)
    ${ }^{1}$ These two authors contributed equally.

[^1]:    ${ }^{2}$ To simplify some of the formulations, we assume that the each $p^{\prime} \neq \mathcal{P}(h)$ takes a noop action at $h$ (which doesn't do anything, and isn't explicitly written into the history, i.e., $h$ noop $=h$ ).

[^2]:    ${ }^{3}$ Intuitively, public states partition the game based on the information that is public - e.g., cards laying face up on the table or actions visible to all players.
    ${ }^{4}$ That is, $(\forall p \in \mathcal{N})\left(\forall I \in \mathcal{I}_{p}\right): I \subset S \vee I \cap S=\emptyset$.
    ${ }^{5}$ Formally: denote by $s_{p}(h)$ the sequence of infosets encountered and actions taken by player $p$ along $h$. Player $p$ has perfect recall if we have $g, h \in I \in \mathcal{I}_{p} \Longrightarrow s_{p}(g)=s_{p}(h)$.
    ${ }^{6}$ Together with $\mathcal{I}_{p}$ covering full $\mathcal{H}$ and players formally taking a noop action in histories where they are inactive, this has the desirable 11 implication of enforcing timeability of all games (informally, the infoset structure being compatible with the players having a sense of time passing (11). For more details, see [7].)

[^3]:    ${ }^{7}$ These results are not particularly difficult to prove and some other researchers might have arrived at them independently. However, we remember not being aware of these results prior to starting the work on this paper, and we consider all of them useful and non-obvious.
    ${ }^{8}$ Where $\sigma(g, a)$ stands for $\sigma_{\mathcal{P}(g)}(I, a)$, where $g \in I \in \mathcal{I}_{\mathcal{P}(g)}$ (resp. $\sigma_{c}(g, a)$ when $\left.\mathcal{P}(g)=c\right)$.

[^4]:    ${ }^{9}$ When $h$ isn't an extension of $g$ (or equal to it), we set $P^{\sigma}(g, h)=0$.
    ${ }^{10}$ Since $v_{p}^{\sigma}(h)$ is the expected utility under $\sigma$, it would also make sense to denote these values simply as $u_{p}^{\sigma}(h)$. The main reason for our preference for $v$ and $q$ is its similarity to the notation used in reinforcement learning, which serves to bring attention to the similarity between the mathematical properties of these objects.

[^5]:    ${ }^{11}$ If we naively defined $P_{p}^{\sigma}(H)$ as a sum $\sum_{h \in H} P_{p}^{\sigma}(h)$, the sum could easily end up strictly larger than 1 , and the resulting quantity would not be useful - e.g., it wouldn't satisfy equation 3.22

[^6]:    ${ }^{12}$ Historically, many implementations of algorithms in EFGs "ran on $\mathcal{H}$ ". The results presented in this section can be particularly relevant when attempting to instead use structures where infosets (or even public states [15]) play a prominent role. Since the infoset trees $\mathcal{I}_{p}$ are much smaller than $\mathcal{H}$, these have the potential to be more effective (particularly if the specific game allows for an efficient computation of $\left.U_{p}^{\sigma}(Z)\right)$.

[^7]:    ${ }^{13}$ Formally, we set $\mathbf{v}\left(h, \sigma^{\mathcal{T}}\right):=0$ when $\operatorname{expl}_{1}\left(\sigma_{1}^{\mathcal{T}}\right)=\operatorname{expl}_{2}\left(\sigma_{2}^{\mathcal{T}}\right)=0, \mathbf{v}\left(h, \sigma^{\mathcal{T}}\right):=-1$ when $\operatorname{expl}_{1}\left(\sigma_{1}^{\mathcal{T}}\right)>0$ and $\operatorname{expl}_{2}\left(\sigma_{2}^{\mathcal{T}}\right)=0, \mathbf{v}\left(h, \sigma^{\mathcal{T}}\right):=1$ when $\operatorname{expl}_{2}\left(\sigma_{2}^{\mathcal{T}}\right)>0$ and $\operatorname{expl} 1\left(\sigma_{1}^{\mathcal{T}}\right)=0$, and $\mathbf{v}\left(h, \sigma^{\mathcal{T}}\right):=0$ when $\operatorname{expl}_{1}\left(\sigma_{1}^{\mathcal{T}}\right), \operatorname{expl}_{2}\left(\sigma_{2}^{\mathcal{T}}\right)>0$.
    ${ }^{14}$ That is, we have $(\forall p \in \mathcal{N})\left(\forall I \in \mathcal{I}_{p}\right): I \cap H \neq \emptyset \Longrightarrow I \subset H$..
    ${ }^{15}$ Recall that since partial strategies are formally subsets of $\mathcal{I}_{p} \times \Delta \mathcal{A}_{p}, \sigma_{p}$ is indeed an extension of $\rho_{p}$ iff $\sigma_{p} \supset \rho_{p}$.

[^8]:    ${ }^{17}$ The metaphor is that if $\mathcal{H}$ is a tree, then $\mathcal{T}$ is its trunk, the subgames rooted at some $S \subset \mathcal{Z}^{\mathcal{T}}$ are its branches, and terminal histories are its leaves. While trunks can and do come in various shapes, we typically imagine trunks of the form $\mathcal{T}=\{h \in \mathcal{H} \backslash \mathcal{Z} \mid$ length $(h) \leq k\}$ for some $k$. (Assuming that game-theoretical trees grow downwards, this is equivalent to $\overline{\mathcal{T}}:=$ "the upper half of $\mathcal{H}$ that has been sliced in two horizontally".)
    ${ }^{18}$ By Definition 4.2 leaves of a trunk $\mathcal{T}$ are just below $\mathcal{T}$. For a trunk strategy profile $\sigma^{\mathcal{T}}$, the leaves of $\mathcal{T}$ are thus the first nodes where $\sigma^{\mathcal{T}}$ becomes undefined.
    ${ }^{18}$ This definition formalizes the notion of trunk used in 17.
    ${ }^{19}$ This is in line with the intuition that a strategy defined nowhere hasn't made any mistakes yet and has, therefore, exploitability zero.
    ${ }^{20}$ Formally, the $G\left(\sigma_{p}^{\mathcal{T}}\right)$ is obtained by turning the decision points in $\sigma_{p}^{\mathcal{T}}$ 's domain into chance nodes with probabilities $\sigma_{p}^{\mathcal{T}}(I)$. Since we can always convert the chance strategies in $\mathcal{T}$ back to players' strategies, we will sometimes abuse the notation slightly and treat strategies in $G\left(\sigma_{p}^{\mathcal{T}}\right)$ as extensions of $\sigma_{p}^{\mathcal{T}}$.

[^9]:    ${ }^{21}$ Recall that "Nash equilibria in the trunk", $\mathrm{NE}(G) \mid \mathcal{T}$, are defined as restrictions of full-game NE to $\mathcal{T}$. To avoid a possible confusion between $\left.\operatorname{NE}(G)\right|_{\mathcal{T}}$ and "Nash equilibria" of $(\mathcal{T}$, $\mathbf{v})$, we will reserve the term solution to depth-limited games and use Nash equilibrium to talk about restrictions of $\mathrm{NE}(G)$.

[^10]:    ${ }^{22}$ Analogously to $G\left(\sigma_{p}^{\mathcal{T}}\right), G\left(\sigma^{\mathcal{T}}\right)$ is obtained by fixing each player's trunk strategy in $G$ to $\sigma^{\mathcal{T}}$. The name value-solving subgame comes from the fact that solving this subgame will result in a strategy with the correct values (however, the strategy might be highly exploitable [17]).
    ${ }^{23}$ The value function $V$ studied by [5] in partially-observable stochastic games is of this type.

[^11]:    ${ }^{24}$ However, such functions do suffice for these purposes if they satisfy additional properties (see Section 4.2.3). One way to satisfy these extra properties is to solve subgames by CFR (in a particular manner), the value functions used in the implementation of $\operatorname{ReBeL}$ [15] are likely to be of this stronger type.

[^12]:    ${ }^{25}$ The different versions of reachability are defined in Section 3.2

[^13]:    ${ }^{26}$ In matching pennies, both players can select either Heads or Tails. If the choices match, player 2 pays one point to player 1. Otherwise, player 1 pays one point to player 2.

[^14]:    ${ }^{27}$ To have the current strategy defined in the whole trunk, including in places that MCTS hasn't yet explored, we can assume that the strategy is initialized to be uniformly random.

[^15]:    ${ }^{28}$ More precisely, by applying CFR to the full game $G$ but keeping the strategies fixed to $\sigma^{t}(I):=\sigma^{\mathcal{T}}(I)$ for all $I \subset \mathcal{T}$ and $t$ (instead of updating them via regret matching as usual).
    ${ }^{29}$ Recall that the D in CFR-D stands for "decomposition" (rather than anything to do with depth-limited solving) 17. Indeed, CFR-D works by decomposing the whole bottom of the game into many smaller subgames. However, this does require first splitting the game into the trunk and the bottom, which makes many of the ideas relevant for depth-limited solving.

[^16]:    ${ }^{30}$ Proposition 4.20 implies that this could be done by fixing the opponent's trunk strategy to $\sigma_{-p}^{\mathcal{T}}$ and solving the resulting 1-player depth-limited game by DL-CFR.

[^17]:    ![](https://cdn.mathpix.com/cropped/2025_03_19_2b8faae0f89ba144f0e1g-32.jpg?height=58&width=1129&top_left_y=2242&top_left_x=482)

[^18]:    ${ }^{32}$ A limitation of Proposition 4.24 is its non-constructiveness - it talks about the solutions of $(\mathcal{T}, \mathbf{v})$ without explaining how to obtain them. Fortunately, if we further assumed that $\mathbf{v}$ is counterfactually optimal (w.r.t. $\mathbb{P}$ ), we could prove that it enables depth-limited variants of constructive methods such as CFR (analogously to how adding counterfactual optimality on top of Theorem 3 gives Proposition 4.20. For brevity, we refrain from discussing this modification in more detail.
    ${ }^{33}$ Recall that, by Definition $4.9 \mathbf{v}$ is said to preserve equilibria of $G$ if any solution of $(\mathcal{T}, \mathbf{v})$ is a restriction of some NE of $G$ to $\mathcal{T}$, and vice versa.
    ${ }^{34}$ By Theorem 4 below, the approach from [4 is equivalent to the one in Proposition 4.24
    ${ }^{35}$ Recall that the simultaneous decision is equivalent to player one choosing first and player two choosing second, without learning the first player's choice.

[^19]:    ${ }^{36}$ While value-functions in other types of POSGSs - e.g., in Dec-POMDPs [23] - are sometimes superficially similar to values in our setting, assumptions often made in these domains (e.g., joint training) give the resulting concepts vastly different properties. As such, their discussion would be outside of the scope of the present paper.
    ${ }^{37}$ Since the notation in this area does not seem to be settled, we should, strictly speaking, say that these concepts are defined on the EFG representation of $G$, or on some object with very similar properties.

[^20]:    ${ }^{38}$ Despite this, the fact remains that using any such value function will lead to good outcomes; the complication lies in the need to be consistent in this choice 47.
    ${ }^{39}$ One might rightfully complain that this counterexample is trivial. However, if optimal values aren't unique in such a simple case, we should expect the situation to be even worse in more complicated games. In particular, we could easily construct a game where player one has ambiguous infoset values even for infosets where they act.

[^21]:    ${ }^{40}$ This is assuming that both players share utility, with player one/two trying to maximize/minimize it.

[^22]:    ${ }^{41}$ A limitation of Proposition 4.24 is its non-constructiveness - it talks about the solutions of $(\mathcal{T}, \mathbf{v})$ without explaining how to obtain them. Fortunately, if we further assumed that $\mathbf{v}$ is counterfactually optimal (w.r.t. $\mathbb{P}$ ), we could prove that it enables depth-limited variants of constructive methods such as CFR (analogously to how adding counterfactual optimality on top of Theorem 3 gives Proposition 4.20 . For brevity, we refrain from discussing this modification in more detail.
    ${ }^{42}$ Recall that, by Definition $4.9 \mathbf{v}$ is said to preserve equilibria of $G$ if any solution of $(\mathcal{T}, \mathbf{v})$ is a restriction of some NE of $G$ to $\mathcal{T}$, and vice versa.

[^23]:    ${ }^{43}$ Recall that for a game $G$ and partial strategy $\rho, G(\rho)$ denotes the game where the players are "forced" to play $\rho$ in infosets where it is defined.

