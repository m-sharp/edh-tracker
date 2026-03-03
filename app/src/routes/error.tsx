import { ReactElement } from "react";
import { useRouteError } from "react-router-dom";

interface RouterError {
    statusText: string;
    message: string;
}

// TODO: This looks bad, content should be flexed and centered on the page
export default function ErrorPage(): ReactElement {
    const error = useRouteError() as RouterError;
    console.error(error);

    return (
        <div id="error-page">
            <h1>Oops!</h1>
            <p>An unexpected error has occurred.</p>
            <p>
                <i>{error.statusText || error.message}</i>
            </p>
        </div>
    );
}
