import { Fragment, PropsWithChildren } from "react";
import { BigTitle } from "../Text/Text";
import styles from "./PageWithSidebar.module.sass";
import { Outlet } from "react-router-dom";

type Props = PropsWithChildren & {
    title: string;
    sidebar: JSX.Element;
};

export default function PageWithSidebar(props: Props) {
    const { title, sidebar, children } = props;
    return (
        <Fragment>
            <BigTitle className={styles.title}>{title}</BigTitle>
            <div className={styles.content}>
                {sidebar}
                <div className={styles.side}>
                    <Outlet />
                </div>
            </div>
        </Fragment>
    );
}
