type Props = {
    name: string;
};

export default function Symbol({ name }: Props) {
    return <span className="material-symbols-rounded">{name}</span>;
}
