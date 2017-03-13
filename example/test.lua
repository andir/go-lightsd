
offset = 0

function render(ctx)
    offset = offset + ctx.duration * 10

    for i = 0, ctx.count - 1 do
        ctx.put((i + offset) % ctx.count,
            i / count * 255,
            i / count * 255,
            i / count * 255)
    end
end
