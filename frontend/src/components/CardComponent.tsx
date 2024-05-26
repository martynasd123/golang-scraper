import React, {HTMLAttributes} from "react";
import "./card.less"
import classNames from "../util/classNames";

const CardComponent: React.FC<HTMLAttributes<HTMLDivElement>> = ({className, ...rest}) => {
    return <div className={classNames("card", className)} {...rest}>
        {rest.children}
    </div>
}

export default CardComponent;