import List from "../../../components/List/List";
import ListItem from "../../../components/List/ListItem";
import ListIcon from "../../../components/List/ListIcon";
import Icon from "../../../components/Icon/Icon";
import ListTitle from "../../../components/List/ListTitle";
import ListInfo from "../../../components/List/ListInfo";
import ListDescription from "../../../components/List/ListDescription";
import { Metric } from "../../../models/metrics";

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
                        <Icon name={icons[metric.type] ?? "show_chart"} />
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
