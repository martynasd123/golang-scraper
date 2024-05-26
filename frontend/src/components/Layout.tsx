import * as React from "react"
import {Outlet} from "react-router-dom";
import "./layout.less"

const Layout: React.FC = () => {
    return <div>
        Layout
        <Outlet/>
    </div>
}

export default Layout