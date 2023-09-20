import { Title } from "../../components/Text/Text";
import Hardware from "../../components/Hardware/Hardware";
import { useFetch } from "../../hooks/useFetch";
import { api } from "../../backend/backend";

import styles from "./SettingsHardware.module.sass";
import { Vertical } from "../../components/Layouts/Layouts";

export default function SettingsHardware() {
    const { data: hardware } = useFetch(api.hardware.get);

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Hardware</Title>
            <Hardware hardware={hardware} />
        </Vertical>
    );
}
