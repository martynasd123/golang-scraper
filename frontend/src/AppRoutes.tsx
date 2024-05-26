import {BrowserRouter, Route, Routes} from "react-router-dom";
import {authenticatedOnly, unauthenticatedOnly} from "./hoc/authProtection";
import Layout from "./components/Layout";
import TasksPage from "./pages/TasksPage";
import TaskProgressPage from "./pages/TaskProgressPage";
import LoginPage from "./pages/LoginPage";
import React from "react";

const AppRoutes: React.FC = () => {
    return <BrowserRouter>
        <Routes>
            <Route path="/" element={authenticatedOnly(Layout)}>
                <Route index element={<TasksPage/>}/>
                <Route path="/task/:id" element={<TaskProgressPage/>}/>
            </Route>
            <Route path="/login">
                <Route index element={unauthenticatedOnly(LoginPage)}/>
            </Route>
        </Routes>
    </BrowserRouter>
}
export default AppRoutes