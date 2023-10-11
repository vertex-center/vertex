import { Title } from "../../../components/Text/Text";
import Hardware from "../../../components/Hardware/Hardware";
import { api } from "../../../backend/backend";

import styles from "./SettingsHardware.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/APIError";
import List from "../../../components/List/List";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery } from "@tanstack/react-query";

export default function SettingsHardware() {
    const {
        data: hardware,
        error,
        isLoading,
    } = useQuery({
        queryKey: ["hardware"],
        queryFn: api.hardware,
    });

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Hardware</Title>
            <APIError error={error} />
            <List>
                <Hardware hardware={hardware} />
            </List>
        </Vertical>
    );
}
