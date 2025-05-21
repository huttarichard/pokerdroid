import { useState } from "react";
import { Box, Stack } from "@mui/material";
import { getTreeState, Action, TreeResponse } from "~/services/tree";
import PathSelector from "./selector";
import StateBox from "./box";
import PokerMatrix from "./matrix";

function TreePage() {
  // Keep track of all nodes in the path
  const [nodes, setNodes] = useState<(TreeResponse | null)[]>([]);
  const [path, setPath] = useState<Action[]>([]);
  const [selectedPathIndex, setSelectedPathIndex] = useState<number>(-1);

  const handleSubmit = async (solutionIndex: number) => {
    try {
      const rootAction: Action = {
        kind: "root",
        solution: solutionIndex,
      };

      const response = await getTreeState([rootAction]);
      setPath([rootAction]);
      setNodes([response]);
      setSelectedPathIndex(0);
    } catch (err) {
      console.error("Failed to load tree:", err);
    }
  };

  const handleNodeClick = async (index: number, actionIdx?: number) => {
    try {
      // If clicking a previous node, truncate path and nodes
      const newPath = path.slice(0, index + 1);
      const newNodes = nodes.slice(0, index + 1);

      if (actionIdx !== undefined) {
        // Add new action and get its state
        const action: Action = {
          kind: "player",
          actionIdx,
        };
        newPath.push(action);

        const response = await getTreeState(newPath);
        newNodes.push(response);
      }

      setPath(newPath);
      setNodes(newNodes);
      setSelectedPathIndex(newPath.length - 1);
    } catch (err) {
      console.error("Failed to get node:", err);
    }
  };

  const handleChanceNode = async (cards: string[]) => {
    try {
      const action: Action = {
        kind: "chance",
        cards,
      };
      const newPath = [...path, action];

      const response = await getTreeState(newPath);
      setPath(newPath);
      setNodes([...nodes, response]);
      setSelectedPathIndex(newPath.length - 1);
    } catch (err) {
      console.error("Failed to process chance node:", err);
    }
  };

  // Helper function to find the selected cards for a chance node
  const getCardsForChanceNode = (index: number): string[] => {
    // If this is a chance node looking at the next action in the path
    if (nodes[index]?.kind === "chance" && index + 1 < path.length) {
      const nextAction = path[index + 1];
      // The next action should have cards if it was a chance node
      if (nextAction.cards) {
        return nextAction.cards;
      }
    }
    return [];
  };

  // Helper function to collect all selected cards from the path
  const getAllSelectedCards = (): string[] => {
    const allCards: string[] = [];

    // Go through the path and collect cards from chance nodes
    path.forEach((action) => {
      if (action.kind === "chance" && action.cards) {
        allCards.push(...action.cards);
      }
    });

    return allCards;
  };

  if (path.length === 0) {
    return <PathSelector onSubmit={handleSubmit} />;
  }

  // Get all selected cards once for the entire render
  const allSelectedCards = getAllSelectedCards();

  return (
    <Stack sx={{ height: "100%", overflow: "hidden" }}>
      <Box
        ref={(el: HTMLElement) => {
          // Auto-scroll to the end when nodes change
          if (el) {
            el.scrollLeft = el.scrollWidth;
          }
        }}
        sx={{
          overflowX: "auto",
          WebkitOverflowScrolling: "touch",
          borderBottom: 1,
          borderColor: "divider",
          // Hide scrollbar while maintaining functionality
          scrollbarWidth: "none", // Firefox
          "&::-webkit-scrollbar": {
            display: "none", // Chrome, Safari, Edge
          },
          msOverflowStyle: "none", // IE and Edge
          cursor: "grab", // Show grab cursor to indicate draggable area
          userSelect: "none", // Prevent text selection during dragging
        }}
        data-lenis-prevent
        // Add mouse events for drag-to-scroll
        onMouseDown={(e) => {
          // Prevent default to avoid text selection
          e.preventDefault();
          // Store initial position
          const el = e.currentTarget;
          const startX = e.pageX;
          const scrollLeft = el.scrollLeft;

          // Change cursor while dragging
          el.style.cursor = "grabbing";

          // Handle mouse move
          const handleMouseMove = (e: MouseEvent) => {
            const x = e.pageX;
            const walk = (x - startX) * 1.5; // Multiply for faster scroll
            el.scrollLeft = scrollLeft - walk;
          };

          // Handle mouse up
          const handleMouseUp = () => {
            document.removeEventListener("mousemove", handleMouseMove);
            document.removeEventListener("mouseup", handleMouseUp);
            el.style.cursor = "grab";
          };

          // Add event listeners
          document.addEventListener("mousemove", handleMouseMove);
          document.addEventListener("mouseup", handleMouseUp);
        }}
      >
        <Stack
          direction="row"
          spacing={0.75}
          sx={{
            p: 1,
            minWidth: "max-content",
          }}
        >
          {path.map((act, index) => {
            // For chance nodes, augment the action with cards if available
            const actionWithCards = { ...act };

            // If this is a chance node, try to find cards from the next action
            if (nodes[index]?.kind === "chance" && !actionWithCards.cards) {
              const selectedCards = getCardsForChanceNode(index);
              if (selectedCards.length > 0) {
                actionWithCards.cards = selectedCards;
              }
            }

            return (
              <StateBox
                key={index}
                node={nodes[index]!}
                action={actionWithCards}
                isSelected={index === selectedPathIndex}
                onClick={() => handleNodeClick(index)}
                onAction={(actionIdx) => handleNodeClick(index, actionIdx)}
                onChanceSelect={handleChanceNode}
                allSelectedCards={allSelectedCards}
              />
            );
          })}
        </Stack>

        <Box
          component="pre"
          sx={{
            ml: 1.5,
            my: 0,
            bgcolor: "transparent",
            borderRadius: 0,
            border: 0,
            fontSize: 9,
            opacity: 0.5,
          }}
        >
          {nodes[selectedPathIndex]?.tree_history}
        </Box>
      </Box>

      <Box
        sx={{
          flex: 1,
          overflow: "auto",
          p: 1,
        }}
        data-lenis-prevent
      >
        <PokerMatrix node={nodes[selectedPathIndex] || { kind: "root" }} />
      </Box>
    </Stack>
  );
}

export default TreePage;
