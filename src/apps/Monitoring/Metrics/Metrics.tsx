import List from "../../../components/List/List";
import ListItem from "../../../components/List/ListItem";
import ListSymbol from "../../../components/List/ListSymbol";
import Symbol from "../../../components/Symbol/Symbol";
import ListTitle from "../../../components/List/ListTitle";
import ListInfo from "../../../components/List/ListInfo";
import ListDescription from "../../../components/List/ListDescription";
import { Metric } from "../../../models/metrics";

type Props = {
    metrics?: Metric[];
};

export default function Metrics(props: Readonly<Props>) {
    if (!props.metrics) return null;

    let symbols = {
        metric_type_on_off: "toggle_on",
        metric_type_number: "show_chart",
    };

    return (
        <List>
            {props.metrics.map((metric) => (
                <ListItem key={metric.name}>
                    <ListSymbol>
                        <Symbol name={symbols[metric.type] ?? "show_chart"} />
                    </ListSymbol>
                    <ListInfo>
                        <ListTitle>{metric.name}</ListTitle>
                        <ListDescription>{metric.description}</ListDescription>
                    </ListInfo>
                </ListItem>
            ))}
        </List>
    );
}
