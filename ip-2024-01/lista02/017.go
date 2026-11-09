package main
import("fmt")

func main (){
	var codigo,vunit,cmenor,cmedio,cmaior,maior,pc,pv,mv,refmv,refmaior,compra,venda float64
	for{
		fmt.Scan(&codigo)
		fmt.Scan(&pc)
		fmt.Scan(&pv)
		fmt.Scan(&vunit)
		compra += pc*vunit
		venda +=pv*vunit
		if vunit>mv{
			mv = vunit
			refmv = codigo
		}
		if (pv*vunit)-(pc*vunit) >maior {
			maior = (pv*vunit)-(pc*vunit)
			refmaior = codigo
		}
		if (pv*vunit)-(pc*vunit) < (pc*vunit)*0.1{
			cmenor ++
		}else if (pv*vunit)-(pc*vunit) < (pc*vunit)*0.2{
			cmedio ++
		}else{
			cmaior++
		}
		fmt.Print("------------------------------------------------------------------------------------------\n")
		fmt.Print("Quantidade de mercadorias que geraram lucro menor que 10% : ",cmenor,"\n")
		fmt.Print("Quantidade de mercadorias que geraram lucro menor que 10% e menor ou igual a 20%: ",cmedio,"\n")
		fmt.Print("Quantidade de mercadorias que geraram lucro maior que 20% : ",cmaior,"\n")
		fmt.Printf("Codigo da mercadoria que gerou maior lucro: %.0f\n", refmaior)
		fmt.Printf("Codigo da mercadoria mais vendida : %.0f\n",refmv)
		fmt.Printf("Valor total de compras: %.2f, valor total de vendas: %.2f e percentual de lucro total: %.2f%%\n",compra,venda,100*(venda-compra)/compra)
		fmt.Print("------------------------------------------------------------------------------------------\n")
	}
}