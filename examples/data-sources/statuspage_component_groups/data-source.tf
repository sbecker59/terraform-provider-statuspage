data "statuspage_component_groups" "default" {
    
    page_id = local.page_id

    filter {
        name = "name"
        values = [ "value_1", "value_2" ]
    }
}