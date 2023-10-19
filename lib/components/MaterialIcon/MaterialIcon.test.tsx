import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import MaterialIcon from "./index.tsx";

test("it renders", () => {
    render(<MaterialIcon icon="deployed_code_update" />);
    const icon = screen.getByText("deployed_code_update");
    expect(icon).toBeInTheDocument();
});

test("it can have a custom class", () => {
    render(
        <MaterialIcon icon="deployed_code_update" className="custom-class" />,
    );
    const icon = screen.getByText("deployed_code_update");
    expect(icon).toHaveClass("custom-class", "material-symbols-rounded");
});
