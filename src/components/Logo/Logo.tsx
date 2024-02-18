import styles from "./Logo.module.sass";

type Props = {
    name?: string;
    iconOnly?: boolean;
};

export default function Logo({ name, iconOnly }: Readonly<Props>) {
    return (
        <div className={styles.logo}>
            <img width={30} src="/images/logo.png" alt="Logo" />
            {!iconOnly && <h1>{name ?? "Vertex"}</h1>}
        </div>
    );
}
