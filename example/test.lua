
function init()
    offset = 0
end

function update(duration)
    offset = offset + duration * 10
end

function render()
    for i = 0, count - 1 do
        put((i + offset) % count,
            i / count * 255,
            i / count * 255,
            i / count * 255)
    end
end
