# Pages cannot be created via the API, so this is the only way to bring one
# under management — the import ID is just the page code, not a composite
# "page_id/id" like statuspage_component and friends use.
terraform import statuspage_page.my_page my_page_id
