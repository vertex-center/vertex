import styles from "./NoItems.module.sass";
import { Paragraph } from "@vertex-center/components";
import { IconContext } from "@phosphor-icons/react";
import React from "react";

type Props = {
    icon?: React.JSX.Element;
    text?: string;
};

export default function NoItems(props: Readonly<Props>) {
    const { icon, text } = props;

    return (
        <div className={styles.card}>
            <div className={styles.icon}>
                <IconContext.Provider value={{ size: 50, weight: "light" }}>
                    {icon}
                </IconContext.Provider>
            </div>
            <Paragraph className={styles.text}>{text}</Paragraph>
        </div>
    );
}
