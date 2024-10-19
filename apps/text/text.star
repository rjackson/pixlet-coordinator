load("render.star", "render")
load("schema.star", "schema")

def main(config):
    text = config.get("text", "teapot")

    return render.Root(
        delay = 2000,
        child = render.Box(
            render.Row(
                expanded=True,
                main_align="space_evenly",
                cross_align="center",
                children = [
                    render.Text(text),
                ]
            )
        )
    )

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Text(
                id = "text",
                name = "Text",
                desc = "",
                icon = "",
            )
        ]
    )