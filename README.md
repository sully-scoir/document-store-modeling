# Document Store Modeling
Explore different ways to store data in a document store. MongoDB is the doucment store used for testing. 

## Setup
* MongoDB must be installed locally and accessibe via `mongodb://localhost:27017`
* gingko (required for aligning output with specs)

## Run
    cd nested_document_array
    ginkgo -v

## Output
    Running Suite: NestedDocumentArray Suite
    ========================================
    Random Seed: 1656455377
    Will run 8 of 8 specs

    Nested Document Array Size of nested array documents
    When size is10 nested documents
    /Users/sully/Development/document_store_modeling/nested_document_array/size_test.go:17
    Nested Doc Count: 10, Avg Object Size: 815 bytes
    •
    ------------------------------
    Nested Document Array Size of nested array documents
    When size is100 nested documents
    /Users/sully/Development/document_store_modeling/nested_document_array/size_test.go:17
    Nested Doc Count: 100, Avg Object Size: 7926 bytes
    •
    ------------------------------
    Nested Document Array Size of nested array documents
    When size is1000 nested documents
    /Users/sully/Development/document_store_modeling/nested_document_array/size_test.go:17
    Nested Doc Count: 1000, Avg Object Size: 80827 bytes
    •
    ------------------------------
    Nested Document Array Size of nested array documents
    When size is10000 nested documents
    /Users/sully/Development/document_store_modeling/nested_document_array/size_test.go:17
    Nested Doc Count: 10000, Avg Object Size: 827828 bytes
    •
    ------------------------------
    Nested Document Array Query by Name
    Client-side query given document Id
    /Users/sully/Development/document_store_modeling/nested_document_array/querying_test.go:23
    Runtime: 29.301278ms
    •
    ------------------------------
    Nested Document Array Query by Name
    Client-side query matching array element
    /Users/sully/Development/document_store_modeling/nested_document_array/querying_test.go:48
    Runtime: 37.514132ms
    •
    ------------------------------
    Nested Document Array Query by Name
    Aggregate returning single matching nested array document
    /Users/sully/Development/document_store_modeling/nested_document_array/querying_test.go:77
    Runtime: 27.898469ms
    •
    ------------------------------
    Nested Document Array Query by Name
    Aggregate returning multiple matching nested array documents
    /Users/sully/Development/document_store_modeling/nested_document_array/querying_test.go:115
    Runtime: 24.933745ms
    •
    Ran 8 of 8 Specs in 0.312 seconds
    SUCCESS! -- 8 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS

    Ginkgo ran 1 suite in 1.273627137s
    Test Suite Passed
