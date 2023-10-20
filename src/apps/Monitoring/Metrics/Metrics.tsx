import List from "../../../components/List/List";
import ListItem from "../../../components/List/ListItem";
import ListIcon from "../../../components/List/ListIcon";
import ListTitle from "../../../components/List/ListTitle";
import ListInfo from "../../../components/List/ListInfo";
import ListDescription from "../../../components/List/ListDescription";
import { Metric } from "../../../models/metrics";
import { MaterialIcon } from "@vertex-center/components";

type Props = {
    metrics?: Metric[];
};

export default function Metrics(props: Readonly<Props>) {
    if (!props.metrics) return null;

    let icons = {
        metric_type_on_off: "toggle_on",
        metric_type_number: "show_chart",
    };

    return (
        <List>
            {props.metrics.map((metric) => (
                <ListItem key={metric.name}>
                    <ListIcon>
                        <MaterialIcon
                            icon={icons[metric.type] ?? "show_chart"}
                        />
                    </ListIcon>
                    <ListInfo>
                        <ListTitle>{metric.name}</ListTitle>
                        <ListDescription>{metric.description}</ListDescription>
                    </ListInfo>
                </ListItem>
            ))}
        </List>
    );
}
