import * as React from "react";
import "./loader.less"
import {HTMLAttributes} from "react";

const LoaderComponent: React.FC<HTMLAttributes<HTMLDivElement>> = (props) => {
    return <div {...props}>
        <div className="loader"></div>
    </div>
}

export default LoaderComponent;