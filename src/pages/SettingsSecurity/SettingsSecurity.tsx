import { Title } from "../../components/Text/Text";

import styles from "./SettingsSecurity.module.sass";
import { Vertical } from "../../components/Layouts/Layouts";
import SSHKey, { SSHKeys } from "../../components/SSHKey/SSHKey";

type Props = {};

export default function SettingsSecurity(props: Props) {
    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Security</Title>
            <SSHKeys>
                <SSHKey />
            </SSHKeys>
        </Vertical>
    );
}
