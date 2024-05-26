import * as React from "react"
import {Outlet} from "react-router-dom";
import "./layout.less"
import HeaderComponent from "./HeaderComponent";
import {HeaderContextProvider} from "../contexts/HeaderContext";

const Layout: React.FC = () => {
    return <HeaderContextProvider>
        <div className="layout">
            <HeaderComponent/>
            <div className="page-content-wrapper">
                <div className="page-content">
                    <Outlet/>
                </div>
            </div>
        </div>
    </HeaderContextProvider>
}

export default Layout