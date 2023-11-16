import styles from "./NoItems.module.sass";
import { MaterialIcon, Paragraph } from "@vertex-center/components";

type Props = {
    icon?: string;
    text?: string;
};

export default function NoItems(props: Readonly<Props>) {
    const { icon, text } = props;

    return (
        <div className={styles.card}>
            <MaterialIcon className={styles.icon} icon={icon} />
            <Paragraph className={styles.text}>{text}</Paragraph>
        </div>
    );
}
