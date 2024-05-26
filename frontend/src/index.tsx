import React from 'react';
import ReactDOM from 'react-dom/client';
import {BrowserRouter, Route, Routes} from "react-router-dom";
import AuthenticatedBoundary from "./pages/AuthenticatedBoundary";
import TasksPage from "./pages/TasksPage";
import TaskProgressPage from "./pages/TaskProgressPage";
import LoginPage from "./pages/LoginPage";
import AuthContextProvider from "./contexts/AuthenticationContext";

ReactDOM.createRoot(document.getElementById('root')).render(
    <AuthContextProvider>
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<AuthenticatedBoundary/>}>
                    <Route index element={<TasksPage/>}/>
                    <Route path="/task/:id" element={<TaskProgressPage/>}/>
                </Route>
                <Route path="/login">
                    <Route index element={<LoginPage/>}/>
                </Route>
            </Routes>
        </BrowserRouter>
    </AuthContextProvider>
);