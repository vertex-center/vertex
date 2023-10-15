import styles from "./NoItems.module.sass";
import Icon from "../Icon/Icon";
import { Text } from "../Text/Text";

type Props = {
    icon?: string;
    text?: string;
};

export default function NoItems(props: Readonly<Props>) {
    const { icon, text } = props;

    return (
        <div className={styles.card}>
            <Icon className={styles.icon} name={icon} />
            <Text className={styles.text}>{text}</Text>
        </div>
    );
}
