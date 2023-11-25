import { Host as HostModel } from "../../../models/hardware";
import { SiApple, SiLinux, SiWindows } from "@icons-pack/react-simple-icons";
import {
    List,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
} from "@vertex-center/components";
import { Fragment } from "react";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";

type HostProps = {
    host?: HostModel;
};

export default function Host(props: Readonly<HostProps>) {
    if (!props.host) return null;

    const {
        os,
        hostname,
        platform,
        platform_version,
        kernel_arch,
        uptime,
        boot_time,
    } = props.host;

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

    const uptimeHours = Math.round((uptime / 3600) * 100) / 100;
    const bootTime = new Date(boot_time * 1000).toLocaleString();

    return (
        <Fragment>
            <List>
                <ListItem>
                    <ListIcon>{icon}</ListIcon>
                    <ListInfo>
                        <ListTitle>{hostname}</ListTitle>
                        <ListDescription>
                            {platform}({platform_version}) - {kernel_arch}
                        </ListDescription>
                    </ListInfo>
                </ListItem>
            </List>
            <KeyValueGroup>
                <KeyValueInfo name="Boot Time" icon="restart_alt">
                    {bootTime}
                </KeyValueInfo>
                <KeyValueInfo name="Uptime" icon="arrow_upward">
                    {uptimeHours} days
                </KeyValueInfo>
            </KeyValueGroup>
        </Fragment>
    );
}
