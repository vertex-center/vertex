import { Fragment, useContext } from "react";
import Button from "../../components/Button/Button";
import { Vertical } from "../../components/Layouts/Layouts";
import { Title } from "../../components/Text/Text";
import { ThemeContext } from "../../index";

type Props = {};

const themes = [
    { key: "vertex-dark", label: "Vertex Dark" },
    { key: "vertex-light", label: "Vertex Light" },
    { key: "catppuccin-mocha", label: "Catppuccin Mocha" },
    { key: "catppuccin-macchiato", label: "Catppuccin Macchiato" },
    { key: "catppuccin-frappe", label: "Catppuccin Frappé" },
    { key: "catppuccin-latte", label: "Catppuccin Latte" },
];

export default function SettingsTheme(props: Props) {
    const { theme, setTheme } = useContext(ThemeContext);

    return (
        <Fragment>
            <Title>Theme</Title>
            <Vertical gap={8}>
                {themes.map((t) => (
                    <Button
                        onClick={() => setTheme(t.key)}
                        selectable
                        selected={theme === t.key}
                    >
                        {t.label}
                    </Button>
                ))}
            </Vertical>
        </Fragment>
    );
}
