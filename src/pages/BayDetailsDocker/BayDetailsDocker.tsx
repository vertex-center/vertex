import { Fragment } from "react";
import { Text, Title } from "../../components/Text/Text";

import styles from "./BayDetailsDocker.module.sass";

export default function BayDetailsDocker() {
    return (
        <Fragment>
            <Title>Docker</Title>
            <Text>
                Docker is <span className={styles.disabled}>disabled</span> for
                this instance.
            </Text>
        </Fragment>
    );
}
