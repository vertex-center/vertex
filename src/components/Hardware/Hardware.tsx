import { Hardware as HardwareModel } from "../../models/hardware";
import {
    SiApple,
    SiDocker,
    SiLinux,
    SiWindows,
} from "@icons-pack/react-simple-icons";
import {
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
} from "@vertex-center/components";
import { Horizontal } from "../Layouts/Layouts";

type Props = {
    hardware?: HardwareModel;
};

export default function Hardware(props: Readonly<Props>) {
    if (!props.hardware) return null;
    if (!props.hardware.host) return null;

    const { dockerized } = props.hardware;
    const { os, arch, platform, version, name } = props.hardware.host;

    let icon = undefined;
    switch (os) {
        case "linux":
            icon = <SiLinux />;
            break;
        case "darwin":
            icon = <SiApple />;
            break;
        case "windows":
            icon = <SiWindows />;
            break;
    }

    const description = (
        <Horizontal alignItems="center" gap={6}>
            {platform}({version}) - {arch}
            {dockerized && <SiDocker size={15} />}
        </Horizontal>
    );

    return (
        <ListItem>
            <ListIcon>{icon}</ListIcon>
            <ListInfo>
                <ListTitle>{name}</ListTitle>
                <ListDescription>{description}</ListDescription>
            </ListInfo>
        </ListItem>
    );
}
