import * as React from "react"
import {Outlet} from "react-router-dom";
import "./layout.less"
import HeaderComponent from "./HeaderComponent";

const Layout: React.FC = () => {
    return <div className="layout">
        <HeaderComponent/>
        <div className="page-content-wrapper">
            <div className="page-content">
                <Outlet/>
            </div>
        </div>
    </div>
}

export default Layout