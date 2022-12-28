'use strict'

export function getPlacesListWithNestedTitles(places) {
    let res = [], map = {}

    places.forEach(place => {
        map[place.id] = {
            'title': place.title,
            'parent_id': place.parent_id,
        }
    });

    places.forEach(place => {
        let title = place.title

        if (map[place.parent_id] !== undefined) {
            title += " (" + map[place.parent_id].title + ")"
        }

        res.push(
            {
                "id": place.id,
                "title": title,
                "parent_id": place.parent_id,
            }
        )
    });

    return res
}