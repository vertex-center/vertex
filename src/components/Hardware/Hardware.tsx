import { Hardware as HardwareModel } from "../../models/hardware";
import { SiApple, SiLinux, SiWindows } from "@icons-pack/react-simple-icons";
import ListItem from "../List/ListItem";
import ListSymbol from "../List/ListSymbol";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListDescription from "../List/ListDescription";

type Props = {
    hardware?: HardwareModel;
};

export default function Hardware(props: Props) {
    if (!props.hardware) return null;
    if (!props.hardware.host) return null;

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

    const description = `${platform} (${version}) - ${arch}`;

    return (
        <ListItem>
            <ListSymbol>{icon}</ListSymbol>
            <ListInfo>
                <ListTitle>{name}</ListTitle>
                <ListDescription>{description}</ListDescription>
            </ListInfo>
        </ListItem>
    );
}
