import Symbol from "../Symbol/Symbol";
import ListItem from "../List/ListItem";
import ListSymbol from "../List/ListSymbol";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListDescription from "../List/ListDescription";
import List, { ListProps } from "../List/List";

export function SSHKeys(props: ListProps) {
    return <List {...props} />;
}

type SSHKeyProps = {
    type: string;
    fingerprint: string;
};

export default function SSHKey(props: SSHKeyProps) {
    const { type, fingerprint } = props;

    return (
        <ListItem>
            <ListSymbol>
                <Symbol name="key" />
            </ListSymbol>
            <ListInfo>
                <ListTitle>SSH Key</ListTitle>
                <ListDescription>
                    {type} - {fingerprint}
                </ListDescription>
            </ListInfo>
        </ListItem>
    );
}
