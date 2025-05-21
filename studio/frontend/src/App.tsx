import { RouterProvider, createHashRouter, redirect } from "react-router-dom";

import Layout from "./layout/layout";

import DashboardPage from "./states/dashboard/route";
import TreePage from "./states/tree/route";
import GamePage from "./states/game/route";

const router = createHashRouter([
  {
    id: "root",
    path: "/",
    // loader() {
    //   // Our root route always provides the user, if logged in
    //   return { user: fakeAuthProvider.username };
    // },
    Component: Layout,
    children: [DashboardPage, TreePage, GamePage],
  },
  {
    path: "/logout",
    async action() {
      // We signout in a "resource route" that we can hit from a fetcher.Form
      return redirect("/");
    },
  },
]);

export default function App() {
  return <RouterProvider router={router} />;
}
