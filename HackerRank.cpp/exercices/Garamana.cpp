#include<bits/stdc++.h>

using namespace std;

string stg1,stg2,teste;
int contador;

int main(){
    ios::sync_with_stdio(false);
    cin.tie(NULL);
    cin>>stg1;
    cin>>stg2;
    int index;
    for (int c=0;c<stg2.size();c++){
        if(stg2[c]=='*'){
            contador++;
        }
    }
    for (int p=0;p<stg1.size();p++){
        for (int s=0;s<stg2.size();s++){
            if (stg1[p]==stg2[s]){
                stg2[s]='0';
                break;
            }else if (contador>0 and s==stg2.size()-1){
                contador--;
                index = stg2.find('*');
                stg2[index]='0';
                break;
            }
        }
    }
    for(int c=0;c<stg1.size();c++){
        teste += '0';
    }
    if (stg2==teste){
        cout<<'S'<<'\n';
    }else{
        cout<<'N'<<"\n";
    }
}
