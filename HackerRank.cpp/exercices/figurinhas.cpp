#include <bits/stdc++.h>
using namespace std;

int n,c,m,temp,contad,flag;
string teste,ptr;
vector<int> posse;
int main(){
    cin.tie(NULL);
    cin>>n;
    cin>>c;
    cin>>m;
    vector<int>carimbadas;
    for(int i=0; i<c;i++){
        cin>>temp;
        carimbadas.push_back(temp);
    }
    for(int x=0;x<m;x++){
        cin>>temp;
        for (int i =0 ;i<c;i++){
            if (temp==carimbadas[i]){
                contad++;
                carimbadas[i]=0;
                break;
            }
        }
    }
    cout<<c-contad<<endl;
}
