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

type SSHKeyProps = {};

export default function SSHKey(props: SSHKeyProps) {
    return (
        <ListItem>
            <ListSymbol>
                <Symbol name="key" />
            </ListSymbol>
            <ListInfo>
                <ListTitle>SSH Key</ListTitle>
                <ListDescription>SSH Key</ListDescription>
            </ListInfo>
        </ListItem>
    );
}
