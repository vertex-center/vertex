import { Hardware as HardwareModel } from "../../models/hardware";
import { SiApple, SiLinux, SiWindows } from "@icons-pack/react-simple-icons";

import styles from "./Hardware.module.sass";
import { Vertical } from "../Layouts/Layouts";

type Props = {
    hardware?: HardwareModel;
};

export default function Hardware(props: Props) {
    if (!props.hardware) return null;

    const { os, arch, version, name } = props.hardware;

    let icon = undefined;
    switch (os) {
        case "linux":
            icon = <SiLinux />;
            break;
        case "Darwin":
            icon = <SiApple />;
            break;
        case "windows":
            icon = <SiWindows />;
            break;
    }

    return (
        <div className={styles.hardware}>
            <div className={styles.icon}>{icon}</div>
            <Vertical gap={4}>
                <div>{name}</div>
                <div className={styles.version}>
                    {version} - {arch}
                </div>
            </Vertical>
        </div>
    );
}
