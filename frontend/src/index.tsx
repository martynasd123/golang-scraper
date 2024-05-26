import React from 'react';
import "./global.less"
import ReactDOM from 'react-dom/client';
import AuthContextProvider from "./contexts/AuthenticationContext";
import {AlertContextProvider} from "./contexts/AlertContext";
import AlertComponent from "./components/AlertComponent";
import AppRoutes from "./AppRoutes";

ReactDOM.createRoot(document.getElementById('root')).render(
    <AuthContextProvider>
        <AlertContextProvider>
            <AlertComponent/>
            <AppRoutes />
        </AlertContextProvider>
    </AuthContextProvider>
);