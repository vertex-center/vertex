import { PropsWithChildren } from "react";
import styles from "./PageWithSidebar.module.sass";
import { Outlet } from "react-router-dom";

type Props = PropsWithChildren & {
    sidebar: JSX.Element;
};

export default function PageWithSidebar(props: Readonly<Props>) {
    const { sidebar } = props;
    return (
        <div className={styles.content}>
            {sidebar}
            <div className={styles.side}>
                <Outlet />
            </div>
        </div>
    );
}
