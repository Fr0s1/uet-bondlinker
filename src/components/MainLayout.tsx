import { Outlet } from "react-router";
import Navbar from "./Navbar";
import ProtectedRoute from "./ProtectedRoute";

export function MainLayout() {
  return (
    <div>
      <Navbar />
      <ProtectedRoute>
        <Outlet />
      </ProtectedRoute>
    </div>
  )
}
