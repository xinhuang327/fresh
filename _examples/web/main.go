package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "345 Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	fmt.Println("https://www.amazon.cn/Webster-La-Aranita-Miedosa-Lucado-Max/dp/0881134791?ie=UTF8&SubscriptionId=AKIAJRVOJOE7DKIAPX3A&camp=2025&creative=165953&creativeASIN=0881134791&linkCode=xm2&tag=pinstame-23")
	fmt.Println("http://www.amazon.cn/Melissa-Doug-%E7%BE%8E%E4%B8%BD%E8%8E%8E%E5%92%8C%E8%B1%86%E8%B1%86-%E5%8F%AF%E9%87%8D%E5%A4%8D%E7%94%A8%E8%B4%B4%E7%BA%B8%E6%9D%BF-%E4%BA%A4%E9%80%9A%E5%B7%A5%E5%85%B7/dp/122308079X?SubscriptionId=AKIAJRVOJOE7DKIAPX3A&tag=pinstame-23&linkCode=xm2&camp=2025&creative=165953&creativeASIN=122308079X")
	fmt.Println("!!Hello from the outside!!~~")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
