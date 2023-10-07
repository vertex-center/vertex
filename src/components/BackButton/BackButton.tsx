import { Link } from "react-router-dom";
import Symbol from "../Symbol/Symbol";
import styles from "./BackButton.module.sass";

type Props = {
    to: string;
};

export default function BackButton(props: Readonly<Props>) {
    const { to } = props;

    return (
        <Link to={to} className={styles.button}>
            <Symbol name="arrow_back" />
        </Link>
    );
}
