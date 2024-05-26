import React, {createContext, FC, PropsWithChildren, useState} from "react";

export enum AlertType {
    WARNING,
    SUCCESS
}

export interface Alert {
    message: string,
    type: AlertType,
    key?: string
}

interface AlertContextT {
    showAlert: (alert: Alert) => void;
    alerts: Alert[]
}

export const AlertContext = createContext<AlertContextT>(null);

export const AlertContextProvider: FC<PropsWithChildren> = ({children}) => {
    const [alerts, setAlerts] = useState<Alert[]>([]);

    const showAlert = (alert: Alert) => {
        const key = Date.now().toString()
        setAlerts(prev => ([...prev, {
            ...alert,
            key
        }]));
        setTimeout(() => setAlerts(prev => prev.filter(alert => alert.key != key)), 3000);
    }

    return <AlertContext.Provider value={{ showAlert, alerts }}>
        {children}
    </AlertContext.Provider>
}