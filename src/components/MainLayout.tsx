import { Outlet } from "react-router";
import Navbar from "./Navbar";
import ProtectedRoute from "./ProtectedRoute";

export function MainLayout() {
  return (
    <div>
      <Navbar />
      <ProtectedRoute>
        <div className="px-4 max-w-screen-2xl m-auto">
          <Outlet />
        </div>
      </ProtectedRoute>
    </div>
  )
}
