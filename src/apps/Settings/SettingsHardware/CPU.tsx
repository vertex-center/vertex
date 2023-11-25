import { CPU as CPUModel } from "../../../models/hardware";
import {
    List,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
} from "@vertex-center/components";

type HostProps = {
    cpu?: CPUModel;
};

export default function CPU(props: Readonly<HostProps>) {
    if (!props.cpu) return null;

    const { model_name, mhz, cores_count } = props.cpu;

    return (
        <List>
            <ListItem>
                <ListIcon>
                    <MaterialIcon icon="memory" />
                </ListIcon>
                <ListInfo>
                    <ListTitle>{model_name}</ListTitle>
                    <ListDescription>
                        {mhz} MHz - {cores_count} cores
                    </ListDescription>
                </ListInfo>
            </ListItem>
        </List>
    );
}
