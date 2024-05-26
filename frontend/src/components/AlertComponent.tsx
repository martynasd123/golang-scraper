import React, {useContext} from "react";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import classNames from "../util/classNames";
import "./alert.less"

const AlertComponent: React.FC = ({}) => {
    const {alerts} = useContext(AlertContext);
    return <div className="alerts-container">
        {alerts.map(alert => (<div
            className={classNames("alert", {
                    "warning": alert.type == AlertType.WARNING,
                    "success": alert.type == AlertType.SUCCESS
                })}
            key={alert.key}>
            {alert.message}
        </div>))}
    </div>
}

export default AlertComponent;