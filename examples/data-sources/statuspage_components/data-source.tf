data "statuspage_components" "default" {
    
    page_id = local.page_id

    filter {
        name = "name"
        values = [ "value_1", "value_2" ]
    }
}