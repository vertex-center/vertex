import { Fragment } from "react";
import { Title } from "../../components/Text/Text";

import styles from "./SettingsHardware.module.sass";
import Hardware from "../../components/Hardware/Hardware";
import { useFetch } from "../../hooks/useFetch";
import { api } from "../../backend/backend";

export default function SettingsHardware() {
    const { data: hardware } = useFetch(api.hardware.get);

    return (
        <Fragment>
            <Title className={styles.title}>Hardware</Title>
            <Hardware hardware={hardware} />
        </Fragment>
    );
}
