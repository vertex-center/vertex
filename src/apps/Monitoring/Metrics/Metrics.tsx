import List from "../../../components/List/List";
import ListItem from "../../../components/List/ListItem";
import ListSymbol from "../../../components/List/ListSymbol";
import Symbol from "../../../components/Symbol/Symbol";
import ListTitle from "../../../components/List/ListTitle";
import ListInfo from "../../../components/List/ListInfo";
import ListDescription from "../../../components/List/ListDescription";

type Props = {
    metrics?: Metric[];
};

export default function Metrics(props: Props) {
    if (!props.metrics) return null;

    return (
        <List>
            {props.metrics.map((metric) => (
                <ListItem key={metric.name}>
                    <ListSymbol>
                        <Symbol name="equalizer" />
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
